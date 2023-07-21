# Bitnet Core
Bitnet is decentralized like Bitcoin, but powerful like Ethereum. Bitnet is a decentralized technology that powers programmable money. Much like Bitcoin, users can make use of the protocol to generate, send, and receive Bitnets among themselves, and like Ethereum, it is also compatible with smart contract and can be used to develop complex subprotocols that can power programmable money, digital assets, non-fungibles, decentralized autonomous organizations and more.

It has no pre-mined supply, and for a while after the launch, anyone with a regular PC will be able to setup a node, mine, and generate Bitnets. By doing that you will help securing the network and keeping it decentralized. Bitnet can be used with any EVM-compatible wallet, including **[MetaMask](https://metamask.io/download/)**.

*There is only one way forward, and it is decentralization.*

**[Visit Website](https://bitnet.money/)**  
**[Wiki/Docs](/BitnetMoney/bitnet/wiki)**
**[Download](https://github.com/BitnetMoney/bitnet/releases)**

## Building from Source

To build from the source, you will need both **[Golang 1.19+](https://go.dev/dl/)** and a **C Compiler** installed to build Bitnet. If any of the two is missing or corrupt, your build will not work.
  
*You don't need to build from the source to run a local node. For that, you can just download the latest pre-built version for your operational system and use it to run your node. Go to the [the releases page](/BitnetMoney/bitnet/releases) for download links.*

### Linux/MacOS
You can build Bitnet in your Linux or MacOS device using the command below:

```
git clone https://github.com/BitnetMoney/bitnet.git && cd bitnet && go run build/ci.go install ./cmd/bitnet
```

### Windows
You can build Bitnet in your Windows device using the command below:
```
git clone https://github.com/BitnetMoney/bitnet.git ; cd bitnet ; .\build.win
```

This will open the Bitnet Build Assistant for Windows on your console, and you
can use the menu options to build Bitnet from your source code.

## Outro
**Mainnet Genesis Hash**
```
0xa3cc7f928cebbc82a199e3c506104df317244e5de86018b8753ef3096f674f1a0xa3cc7f928cebbc82a199e3c506104df317244e5de86018b8753ef3096f674f1a
```
