#!/bin/bash

# This script builds Bitnet from the source code in
# Unix devices. It still requires Goland and a C
# compiler installed to work properly.

# SCRIPT NOT TESTED

echo "Cleaning cache and previous builds..."
    # Cleaning GO cache.
    go clean -cache
    # Looking for old files and deleting them.
    rm -fr build/_workspace/pkg/ ${GOBIN}/*
    rm -r build/bin/bitnet
    rm -r build/bin/abidump
    rm -r build/bin/abigen
    rm -r build/bin/bootnode
    rm -r build/bin/clef
    rm -r build/bin/devp2p
    rm -r build/bin/bitnetkey
    rm -r build/bin/evm
    rm -r build/bin/faucet
    rm -r build/bin/p2psim
    rm -r build/bin/rldpdump

echo "Building Bitnet. This process may take several minues..."
    sudo go run build/ci.go install
    sudo chmod 755 build/bin/geth
    sudo mv build/bin/geth build/bin/bitnet
    sudo chmod 755 build/bin/bitnet

echo "Deploying Bitnet to your local binaries folder..."
    sudo mv build/bin/bitnet /usr/bin

echo "Finished building and deploying Bitnet!"
    echo "You can run it by using the command - bitnet - in your terminal."
    echo ""