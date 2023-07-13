⚠ This software is currently under development and not ready for use. Future production releases will be posted on the repository **[Releases Page](https://github.com/masayoshikob/bitnet/releases)**.

# Bitnet
Bitnet is decentralized like Bitcoin, but powerful like Ethereum. Bitnet is a decentralized technology that powers programmable money. Much like Bitcoin, users can make use of the protocol to generate, send, and receive Bitnets among themselves, and like Ethereum, it is also compatible with smart contract and can be used to develop complex subprotocols that can power programmable money, digital assets, non-fungibles, decentralized autonomous organizations and more.

It has no pre-mined supply, and for a while after the launch, anyone with a regular PC will be able to setup a node, mine, and generate Bitnets. By doing that you will help securing the network and keeping it decentralized. Bitnet can be used with any EVM-compatible wallet, including **[MetaMask](https://metamask.io/download/)**.

*There is only one way forward, and it is decentralization.*

**[Visit Website](https://bitnet.money/)**
**[Join the Conversation](https://bitnet.money/forum)**   
**[Read the Paper](https://bitnet.money/d/bitnet.pdf)**  

## Running a Node
Firstly, download the latest release available for your operational system using the links available on the **[releases page](https://github.com/masayoshikob/bitnet/releases)**.

### Linux and MacOs
(coming soon)

### Windows
1. Extract the contents of the `.zip` file you downloaded inside the directory you want to store your node information. Make sure you have the correct authorisation level to read/write on the folder you're storing your data, otherwise your node instance might not work.

2. You can modify your node parameters by editing the `.config` file.
*Bitnet will start an RPC node by default, but if you want to run a local node without allowing incoming connections, you can replace the `*` and the `0.0.0.0` settings inside the `.config` with `localhost`.*

3. After you have set your node parameters, you can start your node by executing the `Bitnet.cmd` script.

4. With your node running, you can execute the script `BitnetConsole.cmd` to open the Javascript Console so you can control your node.

## Mining
With your node running, you can start mining using the integrated Javascript Console. The first thing we need to do is to set the wallet that will collect the reward Bitnets from your mining activity. Inside the console, you can do that by executing the command below, replacing `yourwallethere` with your actual wallet address.

```miner.setEtherbase('yourwallethere')```

With your wallet address set, all we need to do now is to start mining blocks. For that, you can use the command below.

```miner.start(1)```

*You can replace the number `1` with the number of processor cores you want to use for the mining activity, and you can use `()` to use all available cores.*

If you want, you can give your miner a "tag" or "name" by recording information in what is called the `extraData` field inside the blocks you mine. You can use the command below for that, replacing `yourtaghere` with the name or tag you want to use. We recommend keeping it below 12 characters.

```miner.setExtra('yourtaghere')```

## Building
To build from the source, you will need both **[Golang](https://go.dev/dl/)** and a **C Compiler** installed to build Bitnet. If any of the two is missing or corrupt, your build will not work.
  
*You don't need to build from the source to run a local node. For that, you can just download the latest pre-built version for your operational system and use it to run your node. Go to the [the releases page](https://github.com/masayoshikob/bitnet/releases) for download links.*

### Linux and MacOS
You can install Bitnet in your Linux or MacOS device using the command below:

```
git clone https://github.com/masayoshikob/bitnet.git && cd bitnet && go run build/ci.go install ./cmd/geth
```

### Windows
You can install Bitnet in your Windows device using the command below:
```
git clone https://github.com/masayoshikob/bitnet.git ; cd bitnet ; .\build.win
```

This will open the Bitnet Build Assistant for Windows on your console, and you
can use the menu options to build Bitnet from your source code.

## Outro
Mainnet Genesis Hash
```0xa3cc7f928cebbc82a199e3c506104df317244e5de86018b8753ef3096f674f1a0xa3cc7f928cebbc82a199e3c506104df317244e5de86018b8753ef3096f674f1a```