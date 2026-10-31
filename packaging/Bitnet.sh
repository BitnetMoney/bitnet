#!/bin/sh
# Starts a Bitnet mainnet node with its data in ./bitnet.db, the same way as the
# official package, so an existing bitnet.db folder keeps working. Extra flags
# are passed on to the node, for example: ./Bitnet.sh --http
cd "$(dirname "$0")" || exit 1

echo "Initiating mainnet genesis..."
./bitnet --datadir bitnet.db init .genesis || exit 1

echo "Starting node..."
exec ./bitnet --networkid 210 --config .nodeconfig "$@"
