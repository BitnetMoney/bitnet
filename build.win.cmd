@echo off

:: Bitnet Build Assistant for Windows v.1.0.0
:: This script is a user-friendly way for less experienced
:: users to build Bitnet from the source code in their
:: Windows devices. It still requires Golang and a C
:: compiler installed to work properly.

:build.home
cls
echo Welcome to Bitnet Build Assistant for Windows
echo.
echo Please select one of the options below:
echo 1. Build Bitnet
echo 2. Build Full Source
echo 3. Build DevTools
echo 4. Clean Build Cache
echo 5. Run Linters
echo 6. Run Test
echo.
set /p opt="_"
echo.
if /i %opt% == 1 (goto build.bitnet)
if /i %opt% == 2 (goto build.all)
if /i %opt% == 3 (goto build.devtools)
if /i %opt% == 4 (goto clean.cache)
if /i %opt% == 5 (goto run.linters)
if /i %opt% == 6 (goto run.test) else (
    echo Invalid option! Please try again.
    pause >NULL
    del NULL
    goto build.home
)

:build.bitnet
    echo Building Bitnet...
    echo Searching for old binaries and cleaning old cache...
    del /f /q build\bin\bitnet.exe >NULL
    del /f /q build\bin\abidump.exe >NULL
    del /f /q build\bin\abigen.exe >NULL
    del /f /q build\bin\bootnode.exe >NULL
    del /f /q build\bin\clef.exe >NULL
    del /f /q build\bin\devp2p.exe >NULL
    del /f /q build\bin\bitnetkey.exe >NULL
    del /f /q build\bin\evm.exe > NULL
    del /f /q build\bin\faucet.exe >NULL
    del /f /q build\bin\p2psim.exe >NULL
    del /f /q build\bin\rlpdump.exe >NULL
    go clean -cache
    go run build/ci.go install ./cmd/bitnet
    echo Build finished. Press any key to continue.
    pause > NULL
    del NULL
    exit

:build.all
    echo Building Bitnet (ALL BINARIES)...
    echo Searching for old binaries and cleaning old cache...
    del /f /q build\bin\bitnet.exe >NULL
    del /f /q build\bin\abidump.exe >NULL
    del /f /q build\bin\abigen.exe >NULL
    del /f /q build\bin\bootnode.exe >NULL
    del /f /q build\bin\clef.exe >NULL
    del /f /q build\bin\devp2p.exe >NULL
    del /f /q build\bin\bitnetkey.exe >NULL
    del /f /q build\bin\evm.exe >NULL
    del /f /q build\bin\faucet.exe >NULL
    del /f /q build\bin\p2psim.exe >NULL
    del /f /q build\bin\rlpdump.exe >NULL
    go clean -cache
    go run build/ci.go install
    echo Build finished. Press any key to continue.
    pause > NULL
    del NULL
    exit

:build.devtools
	go install golang.org/x/tools/cmd/stringer@latest
	go install github.com/fjl/gencodec@latest
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install ./cmd/abigen
    echo Build finished. Press any key to continue.
    pause > NULL
    del NULL
    exit

:clean.cache
    echo Searching for old binary files...
    del /f /q build\bin\bitnet.exe >NULL
    del /f /q build\bin\abidump.exe >NULL
    del /f /q build\bin\abigen.exe >NULL
    del /f /q build\bin\bootnode.exe >NULL
    del /f /q build\bin\clef.exe >NULL
    del /f /q build\bin\devp2p.exe >NULL
    del /f /q build\bin\bitnetkey.exe >NULL
    del /f /q build\bin\evm.exe >NULL
    del /f /q build\bin\faucet.exe >NULL
    del /f /q build\bin\p2psim.exe >NULL
    del /f /q build\bin\rldpdump.exe >NULL
    go clean -cache
    rmdir /s /q build/_workspace/pkg
    echo Cache cleanead. Press any key to continue.
    pause > NULL
    del NULL
    goto build.home

:run.linters
    go run build/ci.go lint
    echo Build finished. Press any key to continue.
    pause > NULL
    del NULL
    exit

:run.test
    go run build/ci.go test
    echo Test finished. Press any key to continue.
    pause > NULL
    del NULL
    exit