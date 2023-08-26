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
:: This script is a user-friendly way for less experienced
:: users to build Bitnet from the source code in their
:: Windows devices. It still requires Golang and a C
:: compiler installed to work properly.

:: Bitnet Build Assistant for Windows v.1.0.2

:home
    cls
    echo Welcome to Bitnet Build Assistant for Windows.
    echo.
    echo Please select one of the options below:
    echo 1. Build Bitnet
    echo 2. Build Full Source
    echo 3. Build DevTools
    echo 4. Clean Build Cache
    echo 5. Run Linters
    echo 6. Run Test
    echo.
    set /p userChoice="Your choice: "
    echo.
    if /i %userChoice%==1 goto build_bitnet
    if /i %userChoice%==2 goto build_all
    if /i %userChoice%==3 goto build_devtools
    if /i %userChoice%==4 goto clean_cache
    if /i %userChoice%==5 goto run_linters
    if /i %userChoice%==6 goto run_test
    echo Invalid option! Please try again.
    pause >nul
    goto home

:clean_old_binaries_and_cache
    echo Cleaning old binaries and cache...
    rmdir /s /q build\bin\
    go clean -cache

:build_bitnet
    echo Building Bitnet...
    call :clean_old_binaries_and_cache
    go run build/ci.go install ./cmd/bitnet
    echo Build finished.
    pause >nul
    exit

:build_all
    echo Building all binaries...
    call :clean_old_binaries_and_cache
    go run build/ci.go install
    echo Build finished.
    pause >nul
    exit

:build_devtools
    echo Installing devtools...
    go install golang.org/x/tools/cmd/stringer@latest
    go install github.com/fjl/gencodec@latest
    go install github.com/golang/protobuf/protoc-gen-go@latest
    go install ./cmd/abigen
    echo Installation finished.
    pause >nul
    exit

:clean_cache
    echo Cleaning cache...
    call :clean_old_binaries_and_cache
    echo Cache cleared.
    pause >nul
    goto home

:run_linters
    echo Running linters...
    go run build/ci.go lint
    echo Linters have been run.
    pause >nul
    exit

:run_test
    echo Running tests...
    go run build/ci.go test
    echo Tests completed.
    pause >nul
    exit
