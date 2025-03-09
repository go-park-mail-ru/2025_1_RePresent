import subprocess
import sys
import os


def install_pip():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "--upgrade", "pip"])


def install_pip_tools():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "pip-tools"])


def generate_requirements():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "pipreqs"])

    try:
        subprocess.check_call(["pipreqs", ".", "--force", "--encoding=utf-8"])
        subprocess.check_call(
            ["pip-compile", "requirements.txt", "-o", "requirements-compiled.txt"]
        )
    except subprocess.CalledProcessError as e:
        print(f"Произошла ошибка при генерации файла требований: {e}")
        exit(1)


def install_requirements():
    requirements_file = "./requirements-compiled.txt"
    if not os.path.exists(requirements_file):
        print(f"Файл {requirements_file}, необходимо вызвать generate_requirements().")
        return

    try:
        subprocess.check_call(
            [sys.executable, "-m", "pip", "install", "-r", requirements_file]
        )
    except subprocess.CalledProcessError as e:
        print(f"Произошла ошибка при установке требований: {e}")
        exit(1)


def configurate():
    install_pip()
    install_pip_tools()
    generate_requirements()
    install_requirements()
