@echo off
:build.home
cls
echo Welcome to Bitnet Build for Windows
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
    pause > NULL
    del NULL
    goto build.home
)

:build.bitnet
    echo Building Bitnet...
    echo Searching for old binaries and cleaning old cache...
    del /f /q build\bin\bitnet.exe
    go clean -cache
    go run build/ci.go install ./cmd/geth
    ren build\bin\geth.exe bitnet.exe
    echo Build finished. Press any key to continue.
    pause > NULL
    del NULL
    exit

:build.all
    echo Building Bitnet (ALL BINARIES)...
    echo Searching for old binaries and cleaning old cache...
    del /f /q build\bin\bitnet.exe
    go clean -cache
    go run build/ci.go install
    ren build\bin\geth.exe bitnet.exe
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