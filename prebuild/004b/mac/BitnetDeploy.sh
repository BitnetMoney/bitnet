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
# This script will deploy 'bitnet' to your '/usr/bin' folder so
# you can use the bitnet cli.

# Deploys bitnet to /usr/bin
echo "Deploying Bitnet..."
sudo mv bitnet /usr/bin
echo "Deployment finished! You can now use the bitnet binary from anywhere."
echo "Don't forget to make the appropriate changes to your script files."