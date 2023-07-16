@echo off

:: Copyright 2023 Bitnet
:: This file is part of the Bitnet library.
::
:: This software is provided "as is", without warranty of any kind,
:: express or implied, including but not limited to the warranties
:: of merchantability, fitness for a particular purpose and
:: noninfringement. In no even shall the authors or copyright
:: holders be liable for any claim, damages, or other liability,
:: whether in an action of contract, tort or otherwise, arising
:: from, out of or in connection with the software or the use or
:: other dealings in the software.
::
:: Bitnet Build Assistant for Windows v.1.0.1
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
         rmdir /s /q build\bin\
            go clean -cache
            go run build/ci.go install ./cmd/bitnet
        echo Build finished. Press any key to continue.
            pause > NULL
            del NULL
            exit

:build.all
    echo Building Bitnet (ALL BINARIES)...
        echo Searching for old binaries and cleaning old cache...
            rmdir /s /q build\bin\
            go clean -cache
            go run build/ci.go install
        echo Build finished. Press any key to continue.
            pause > NULL
            del NULL
            exit

:build.devtools
    echo Installing devtools...
	    go install golang.org/x/tools/cmd/stringer@latest
	    go install github.com/fjl/gencodec@latest
	    go install github.com/golang/protobuf/protoc-gen-go@latest
	    go install ./cmd/abigen
    echo Installation finished. Press any key to continue.
        pause > NULL
        del NULL
        exit

:clean.cache
    echo Searching for old binary files...
        rmdir /s /q build\bin\
        rmdir /s /q build/_workspace/pkg
        go clean -cache
    echo Cache cleanead. Press any key to continue.
        pause > NULL
        del NULL
        goto build.home

:run.linters
    echo Building linters...
        go run build/ci.go lint
    echo Build finished. Press any key to continue.
        pause > NULL
        del NULL
        exit

:run.test
    echo Initiating...
        go run build/ci.go test
    echo Test finished. Press any key to continue.
        pause > NULL
        del NULL
        exit