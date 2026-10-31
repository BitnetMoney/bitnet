@echo off
:: Opens the JavaScript console of the node started with Bitnet.cmd. Extra flags
:: are passed on, for example: BitnetConsole.cmd --exec "eth.blockNumber"
title Bitnet Javascript Console
cd /d "%~dp0"
.\bitnet attach %* \\.\pipe\bitnet
