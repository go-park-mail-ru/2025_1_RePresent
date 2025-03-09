import subprocess
import sys
import os


def install_pip():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "--upgrade", "pip"])


def install_pip_tools():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "pip-tools"])


def install_imports():
    subprocess.check_call(
        [sys.executable, "-m", "pip", "install", "psycopg2-binary", "PyYAML"]
    )


def configurate():
    install_pip()
    install_pip_tools()
    install_imports()
