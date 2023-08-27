# Version History

## v.0.0.4b

- Updated miner binary for Windows
- Consensus improvements
- Security patch to prevent consensus-level attacks
- Updated scripts for all OSs
- General bug fixes
- Parameters updates, with new bootnodes and built-in genesis
- Light updates on the JSConsole

## v.0.0.3b

First beta version of Bitnet.

- Several bug fixes
- Script updates for all operational systems
- Bitnet doesn't start on a RPC node by default anymore
- Security updates with initialization flags
- Updated default port to 30210 to address conflicts with the 30303 port
- Increased default maximum node peers from 50 to 75
- Enabled node discovery by default
- Updated list of static nodes
- Increased default timeout for miner to 30 minutes

## v0.0.2a

Quick version update enabling quick and easy GPU mining, as well as more static node addresses by default for more decentralization.

- Integrated `ethminer` for GPU mining
- New scripts for running GPU mining automatically
- Removed `sudo` from Linux/MacOS scripts
- Linux/MacOS scripts now execute locally
- Added new list of static nodes to the `.config` file

## v0.0.1a

- First release.