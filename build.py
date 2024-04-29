import subprocess
import sys
import os

def handle_error():
    print("Error: Script failed to build Bitnet.")
    sys.exit(1)

def clean_cache():
    print("Cleaning build cache...")
    subprocess.run(["go", "clean", "-cache"], check=True)

def build_bitnet():
    print("Starting to build Bitnet...")
    clean_cache()
    result = subprocess.run(["go", "run", "build/ci.go", "install", "./cmd/bitnet"], check=False)
    if result.returncode != 0:
        handle_error()
    print("Build finished successfully.")

if __name__ == "__main__":
    build_bitnet()
