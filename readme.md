# Bitnet Core: Your Gateway to Decentralized, Programmable Money

Bitnet combines the decentralization of Bitcoin with the programmability of Ethereum, making it a dual-purpose blockchain technology. Utilize the protocol to create, transfer, and receive Bitnets, or leverage its smart contract capabilities to build complex applications—from digital assets and non-fungible tokens (NFTs) to decentralized autonomous organizations (DAOs) and much more.

> **Insight**: *The only way forward is through decentralization.*

## Important Notice For v0.0.4b Update

If you have any issues syncing with the network, please try the following steps:

1. If you're executing the Bitnet binary directly and in a single step, try initializing the protocol genesis before starting the `bitnet` (or `bitnet.exe` on Windows) binary. You can do that by running the command below inside your Bitnet folder - it assumes you have downloaded the latest pre-build release:

```bash
./bitnet --datadir bitnet.db init .bitnet && ./bitnet --networkid 210 --config .config
```

or, if you're on Windows:

```cmd
.\bitnet --datadir bitnet.db init .bitnet ; .\bitnet --networkid 210 --config .config
```

2. Try deleting your existing database for a full resync with the network. You can use the command below to delete the `bitnet.db` folder:

```bash
rm -r bitnet.db
```

or, if you're on Windows:

```cmd
rm bitnet.db
```

If after following the steps above you still cannot sync your node, please ask for help in one of our community channels. Most of our developers are more active on [Discord](https://discord.com/invite/dtw7rKQfRs) than in other social media platforms.

## Build From Source: Pre-Requisites

Should you wish to build Bitnet from source, ensure you have the following software:

- [Golang 1.19+](https://go.dev/dl/)
- A C Compiler

> **Note**: If either of these components is missing or corrupt, the build will fail. For a pre-built version compatible with your operating system, visit our [Releases Page](https://github.com/BitnetMoney/bitnet/releases/).

---

### Build on Linux/MacOS

Execute the following command to build Bitnet on a Linux or MacOS machine:

```bash
git clone https://github.com/BitnetMoney/bitnet.git && cd bitnet && bash build.linux.sh
```

> **MacOS Users**: Substitute `bash build.linux.sh` with `bash build.mac.sh`.

---

### Build on Windows

Run the command below to build Bitnet on a Windows machine:

```cmd
git clone https://github.com/BitnetMoney/bitnet.git ; cd bitnet ; .\build.win
```

This will initiate the Bitnet Build Assistant for Windows, guiding you through the build process.

## Key Metrics

- **Block Reward**: 1
- **Uncle Reward**: 0.875 (1st), 0.75 (2nd), 0.25 (3rd)
- **Pre-mined Coins**: None
- **Total Supply**: ∞
- **Consensus Mechanism**: Proof of Work
- **Algorithm**: Ethash
- **Target Block Time**: 15 seconds
- **Smart Contracts**: Supported

## Genesis Block

The Mainnet Genesis Hash is as follows:

```
0xa3cc7f928cebbc82a199e3c506104df317244e5de86018b8753ef3096f674f1a
```
