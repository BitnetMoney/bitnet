# Bitnet Core
Bitnet is decentralized like Bitcoin, but powerful like Ethereum. Bitnet is a decentralized technology that powers programmable money. Much like Bitcoin, users can make use of the protocol to generate, send, and receive Bitnets among themselves, and like Ethereum, it is also compatible with smart contract and can be used to develop complex subprotocols that can power programmable money, digital assets, non-fungibles, decentralized autonomous organizations and more.

It has no pre-mined supply, and for a while after the launch, anyone with a regular PC will be able to setup a node, mine, and generate Bitnets. By doing that you will help securing the network and keeping it decentralized. Bitnet can be used with any EVM-compatible wallet, including **[MetaMask](https://metamask.io/download/)**.

*There is only one way forward, and it is decentralization.*

**[Visit Website](https://bitnet.money/)**  
**[Join the Conversation](https://bitnet.money/forum)**   
**[Read the Paper](https://bitnet.money/d/bitnet.pdf)**  
  
**[Download](https://github.com/masayoshikob/bitnet/releases)**

# Explore:

**[Running a Node](#running-a-node)**  
**[Generating Bitnets via Mining](#generating-bitnets-via-mining)**  
**[Using MetaMask & Other Wallet Providers](#using-metamask--other-wallet-providers)**  
**[Building from Source](#building-from-source)**  
**[Outro](##building-from-source)**

## Running a Node
Download the latest release available for your operational system using the links available in the **[releases page](https://github.com/masayoshikob/bitnet/releases)**. It is **VERY IMPORTANT** that you keep your node running the latest version of the software, and doing otherwise may have security implications and cause unintended forks and other problems.

Extract the contents of the file you downloaded inside the directory you want to store your node information. Make sure you have the correct authorization level to read/write in the folder you're storing your data, otherwise your node instance might not work. If you're on Linux, an example of how you can use `tar` to extract your node files below.

```
tar –xvzf Unix_Bitnet_v.X.X.X.tar.gz
```
*Replace `Unix_Bitnet_v.X.X.X.tar.gz` with the correct filename before executing the command.*

You can modify your node parameters by editing the `.config` file. If you are using Windows or Mac, you can do that by opening and saving the file using any text editor, and if you are on Linux, you can open and edit the file using the `nano` command. Use `CTRL+D` + `ENTER` to save your modifications, and `CTRL+X` to exit `nano`.
```
nano .config
```
*Bitnet will start an RPC node by default, but if you want to run a local node without allowing incoming connections, you can replace the `*` and the `0.0.0.0` settings inside `.config` with `localhost`.*

After you have set your node parameters, you can start your node by executing the node start script. Examples below:
  
**On Windows Devices:**
```
.\Bitnet
```
**On Unix Devices:***
```
bash bitnet.node.sh
```

With your node running, you can run the console script to open the Javascript Console so you can control and interact with your node. Examples below:
  
**On Windows Devices:**
```
.\BitnetConsole
```
**On Unix Devices:**
```
bash bitnet.console.sh
```

## Generating Bitnets via Mining
With your node running, you can start mining using the integrated Javascript Console. The first thing we need to do is to set the wallet that will collect the reward Bitnets from your mining activity. Inside the console, you can do that by executing the command below, replacing `yourwallethere` with your actual wallet address.

```
miner.setEtherbase('yourwallethere')
```

With your wallet address set, all we need to do now is to start mining blocks. For that, you can use the command below.

```
miner.start(1)
```

*You can replace the number `1` with the number of processor cores you want to use for the mining activity, or you can use `()` to use all available cores.*

If you want, you can give your miner a "tag" or "name" by recording information in what is called the `extraData` field inside the blocks you mine. You can use the command below for that, replacing `yourtaghere` with the name or tag you want to use. We recommend keeping it below 12 characters.

```
miner.setExtra('yourtaghere')
```

## Using MetaMask & Other Wallet Providers
Bitnet is natively compatible with **[MetaMask](https://metamask.io/download/)** and a range of other well established wallet providers. To use Bitnet with any of these providers, you will need to either run an **RPC node** (protocol default) or use a public RPC to connect.
  
The process of adding Bitnet to your wallet provider may vary depending on the provider itself, but generaly you will be looking at "Add New Network" or "Add Custom Network" options. Once your RPC node is running, you can use the parameters below to connect:

- **Network Name:** Bitnet
- **Network ID:** 210
- **RPC URL:** http://127.0.0.1:8545/
- **Currency Symbol:** BTN

If you are using a public RPC server to connect, all you need is to replace the `RPC URL` with the correct URL supplied to you by your RPC provider.

## Building from Source
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
**Mainnet Genesis Hash**
```
0xa3cc7f928cebbc82a199e3c506104df317244e5de86018b8753ef3096f674f1a0xa3cc7f928cebbc82a199e3c506104df317244e5de86018b8753ef3096f674f1a
```