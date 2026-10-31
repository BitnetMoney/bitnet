#!/bin/sh
# Opens the JavaScript console of the node started with Bitnet.sh. Extra flags
# are passed on, for example: ./BitnetConsole.sh --exec "eth.blockNumber"
cd "$(dirname "$0")" || exit 1
exec ./bitnet attach "$@" bitnet.db/bitnet
