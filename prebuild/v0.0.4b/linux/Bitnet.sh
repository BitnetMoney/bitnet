#!/bin/bash

# Copyright 2023 Bitnet
# This file is part of the Bitnet library.
#
# This software is provided "as is", without warranty of any kind,
# express or implied, including but not limited to the warranties
# of merchantability, fitness for a particular purpose and
# noninfringement. In no even shall the authors or copyright
# holders be liable for any claim, damages, or other liability,
# whether in an action of contract, tort or otherwise, arising
# from, out of or in connection with the software or the use or
# other dealings in the software.
#
# This script requires chmod 755.

# Register and initiates the mainnet genesis file.
echo "Initiating mainnet genesis..."
bitnet --datadir bitnet.db init .bitnet

# Starts the node using the parameters specified.
echo "Starting node..."
bitnet --networkid 210 --config .config --miner.etherbase 0x0000000000000000000000000000000000000000 --mine