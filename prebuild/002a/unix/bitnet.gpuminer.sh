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
# This script runs ethermine to enable GPU mining in the pool.
# You can change the software flags below, and for more information
# visit https://github.com/ethereum-mining/ethminer.

./ethminer -P http://127.0.0.1:8545/ --noeval --report-hashrate --response-timeout 60 --work-timeout 300