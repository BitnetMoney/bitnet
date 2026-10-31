# Reproducible builds

A release binary is only as trustworthy as the way it was made. The Linux
binaries of this client are reproducible: anyone can rebuild one from the
tagged source and get a file with exactly the same SHA-256 checksum as the one
published in the release's `SHA256SUMS`. If the checksums match, the binary
contains that source and nothing else.

## What makes them reproducible

- The [`Dockerfile`](../Dockerfile) builds in a Go image pinned by digest. That
  image already contains the C compiler cgo uses, so the digest fixes the whole
  toolchain, not just Go.
- `-trimpath` removes the paths of the machine that built it.
- `-buildvcs=false` leaves out version control state; the commit and its date
  are passed in explicitly, with the date in UTC. A build from a source archive
  and one from a git checkout of the same commit are therefore identical.
- The binary is linked statically, so it runs on any Linux distribution.

The CI workflow builds every change twice from scratch, on amd64 and on arm64,
and fails if the two binaries differ.

## Checking a release yourself

You need git and Docker. Replace `v0.1.0` with the release you are checking.

```bash
git clone https://github.com/<owner>/bitnet && cd bitnet
git checkout v0.1.0
docker build \
  --build-arg COMMIT=$(git rev-parse HEAD) \
  --build-arg DATE=$(TZ=UTC git show -s --format=%cd --date=format-local:%Y%m%d HEAD) \
  --target binary --output out .
sha256sum out/bitnet
```

Compare the result with the `bitnet-linux-amd64` (or `bitnet-linux-arm64`)
line of the release's `SHA256SUMS`. Build on a machine of the same
architecture as the binary you are checking.

## macOS and Windows

The macOS and Windows binaries are built on GitHub's runners with
`go run build/ci.go install` and Go 1.27.3. They use the same flags, so two
builds with the same Go, Xcode or MinGW versions are identical, but those
system compilers are not pinned the way the Linux image is. Check the Linux
binaries when you need a guarantee.

## Running the image

```bash
docker run -d --name bitnet -v bitnet-data:/data -p 30210:30210 -p 30210:30210/udp \
  ghcr.io/<owner>/bitnet:v0.1.0
```

The node keeps its data in the `/data` volume and runs as an unprivileged user.
It starts with the same `.nodeconfig` as the release packages (port 30210 and
the configured peers); flags after the image name are passed on to the node.
Add `--http --http.addr 0.0.0.0` and publish port 8545 only if you need RPC
from outside the container, and expose only the modules you need.
