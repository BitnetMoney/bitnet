@echo off
:: Starts a Bitnet mainnet node with its data in .\bitnet.db, the same way as the
:: official package, so an existing bitnet.db folder keeps working. Extra flags
:: are passed on to the node, for example: Bitnet.cmd --http
title Bitnet Node
cd /d "%~dp0"

echo Initiating mainnet genesis...
.\bitnet --datadir bitnet.db init .genesis
if errorlevel 1 (
  pause
  exit /b 1
)

echo Starting node...
.\bitnet --networkid 210 --config .nodeconfig %*
