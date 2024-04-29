# This script installs the necessary dependencies to build Bitnet, such as GCC and Golang
# Administrative access is required to run this script
# This script is compatible with Windows, Linux, and macOS

import sys
import subprocess

def run_command(command):
    try:
        subprocess.run(command, check=True)
        return True
    except subprocess.CalledProcessError:
        return False

def install_homebrew():
    print("Installing Homebrew...")
    return run_command('/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"')

def install_chocolatey():
    print("Installing Chocolatey...")
    return run_command(['powershell', 'Set-ExecutionPolicy', 'Bypass', '-Scope', 'Process', '-Force', ';',
                        'iwr', 'https://chocolatey.org/install.ps1', '-UseBasicParsing', '|', 'iex'])

def check_or_install_package_manager():
    if sys.platform.startswith('darwin'):
        # Check if Homebrew is installed
        if not run_command(['which', 'brew']):
            if not install_homebrew():
                print("Failed to install Homebrew.")
                sys.exit(1)
        return 'brew'
    elif sys.platform.startswith('win'):
        # Check if Chocolatey is installed
        if not run_command(['choco', '--version']):
            if not install_chocolatey():
                print("Failed to install Chocolatey.")
                sys.exit(1)
        return 'choco'
    return None

def install_package(package, installer):
    print(f"Attempting to install {package}...")
    if not run_command([installer, 'install', package]):
        print(f"Failed to install {package}. Please check your installation settings and permissions.")

def main():
    if sys.platform.startswith('linux'):
        # Assuming Linux users have sudo privileges and basic package managers like apt or yum
        # Extend or modify as per the targeted Linux distributions
        install_package('golang', 'sudo apt')
        install_package('gcc', 'sudo apt')
    elif sys.platform.startswith('darwin'):
        installer = check_or_install_package_manager()
        install_package('go', installer)
        install_package('gcc', installer)
    elif sys.platform.startswith('win'):
        installer = check_or_install_package_manager()
        install_package('golang', installer)
        install_package('mingw', installer)  # mingw includes gcc
    else:
        print("Unsupported operating system.")
        sys.exit(1)

if __name__ == "__main__":
    main()
