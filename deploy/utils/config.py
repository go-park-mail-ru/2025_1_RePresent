import subprocess
import sys
import os

def install_pip():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "--upgrade", "pip"])
    
def install_pip_tools():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "pip-tools"])

def generate_requirements():
    subprocess.check_call([sys.executable, "-m", "pip", "install", "pipreqs"])
    subprocess.check_call(["pipreqs", ".", "--force"])
    subprocess.check_call(["pip-compile", "requirements.txt", "-o", "requirements-compiled.txt"])

def install_requirements():
    requirements_file = "./requirements-compiled.txt"
    if not os.path.exists(requirements_file):
        print(f"Файл {requirements_file} не существует")
        return
    subprocess.check_call([sys.executable, "-m", "pip", "install", "-r", requirements_file])
