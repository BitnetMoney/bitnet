# Running a public RPC endpoint

Wallets, explorers and exchanges need nodes that answer JSON-RPC. This guide is
for anyone who wants to offer one: it keeps the node itself out of reach and
puts a small, well-understood layer in front of it. It needs no third-party
service.

## What the client already enforces

- The `admin`, `debug`, `personal` and `miner` APIs control the node or reveal
  its internals. The client refuses to start if they are served over HTTP or
  WebSocket on an address other machines can reach. `--rpc.allow-unsafe-public`
  overrides this; do not use it on a public endpoint.
- A JSON-RPC batch may contain at most 1,000 calls
  (`--rpc.batch-request-limit`) and produce at most 25 MB of responses
  (`--rpc.batch-response-max-size`). Lower them for a busy public endpoint.
- `eth_call` and `eth_estimateGas` are capped by `--rpc.gascap` (gas) and
  `--rpc.evmtimeout` (time).

## Layout

```
internet ──> 443/tcp  nginx (TLS, POST / only, rate limits) ──> 127.0.0.1:8545  bitnet
internet <─> 30210/tcp+udp  bitnet p2p
```

The node listens for RPC on 127.0.0.1 only. nginx can run on the same machine;
if it runs on another one, bind the node to a private address and keep the
unsafe APIs off as above.

## The node

Run it as its own unprivileged user, for example with systemd:

```ini
# /etc/systemd/system/bitnet.service
[Unit]
Description=Bitnet node
After=network-online.target

[Service]
User=bitnet
ExecStart=/usr/local/bin/bitnet --datadir /var/lib/bitnet \
  --http --http.addr 127.0.0.1 --http.port 8545 --http.vhosts rpc.example.org \
  --http.api eth,net,web3 \
  --rpc.batch-request-limit 100 --rpc.gascap 50000000 --rpc.evmtimeout 5s
Restart=on-failure
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
ReadWritePaths=/var/lib/bitnet

[Install]
WantedBy=multi-user.target
```

Add `bitnet` to `--http.api` if you want to serve the supply, mining share and
confirmation figures; it is read-only.

## nginx

Only `POST /` reaches the node; every other path is answered by nginx. Requests
for `/wp-admin`, `.php` files and the like never get further than a 404.

```nginx
limit_req_zone  $binary_remote_addr zone=rpc:10m rate=20r/s;
limit_conn_zone $binary_remote_addr zone=rpcconn:10m;

server {
    listen 443 ssl;
    server_name rpc.example.org;
    # ssl_certificate ... (for example from Let's Encrypt)

    client_max_body_size 256k;

    location = / {
        limit_except POST { deny all; }
        limit_req  zone=rpc burst=40 nodelay;
        limit_conn rpcconn 20;
        proxy_pass http://127.0.0.1:8545;
        proxy_set_header Host rpc.example.org;
        proxy_read_timeout 30s;
    }

    location / {
        return 404;
    }
}
```

For WebSocket, add a second `location` that proxies to the node's `--ws` port
with the `Upgrade` and `Connection` headers set.

## Fail2Ban

Ban addresses that keep probing paths that do not exist or keep hitting the
rate limit:

```ini
# /etc/fail2ban/filter.d/bitnet-rpc.conf
[Definition]
failregex = ^<HOST> .* "(GET|POST|HEAD|PUT)[^"]*" (404|405|429) 
ignoreregex =

# /etc/fail2ban/jail.d/bitnet-rpc.conf
[bitnet-rpc]
enabled  = true
port     = https
filter   = bitnet-rpc
logpath  = /var/log/nginx/access.log
findtime = 60
maxretry = 30
bantime  = 3600
```

## Firewall

```bash
ufw default deny incoming
ufw allow 443/tcp
ufw allow 30210/tcp
ufw allow 30210/udp
ufw allow OpenSSH
ufw enable
```

Never open the node's RPC port (8545) itself. Opening UDP 30210 as well as TCP
lets other nodes find yours through discovery, so it can serve as a bootnode.

## A content delivery network in front

A CDN or web application firewall can absorb large floods, but it also makes a
single company a gatekeeper for access to the chain. It is optional: the setup
above already rejects what such rules would block, and the limits in the client
apply either way.
