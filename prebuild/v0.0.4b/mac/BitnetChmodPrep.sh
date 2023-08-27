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
# This script gives Chmod 755 to all required binaries and scripts.

echo "Initializing Chmod script..."
sudo chmod 755 bitnet
sudo chmod 755 Bitnet.sh
sudo chmod 755 BitnetConsole.sh
sudo chmod 755 BitnetGPUMiner.sh
echo "All files processed! Now you can start your node with all"
echo "the required permissions."