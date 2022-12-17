import platform
import os
import subprocess


print("=====================")
print("| FLUTTER INSTALLER |")
print("=====================")
print(" ")
print("System Specifications: ")
print("-----------------------")

# System Check Variables
cpu = platform.processor()
os = platform.system()

print("CPU: ", cpu)
print("OS: ", os)

# Check Application Requirements
print(" ")
print("Application Requirements: ")
print("--------------------------")

# Xcode CLI
try:
    subprocess.run(["xcode-select", "--version"], check=True)
    xCLI = ("Installed")
except subprocess.CalledProcessError:
    xCLI = ("Not Installed")

print("Xcode Command Line Tools: ", xCLI)

# Homebrew


def is_brew_installed():
    try:
        subprocess.run(['which', 'brew'], check=True, stdout=subprocess.PIPE)
        return True
    except subprocess.CalledProcessError:
        return False


if is_brew_installed():
    brew = ("Installed.")
else:
   brew = ("Not Installed.")


print("Homebrew: ", brew)

# Cocoapods


def is_cocoapods_installed():
    try:
        subprocess.run(['which', 'pod'], check=True, stdout=subprocess.PIPE)
        return True
    except subprocess.CalledProcessError:
        return False


if is_cocoapods_installed():
    cPods = ("Installed")
else:
    cPods = ("Not Installed")

print("Cocoapods: ", cPods)


# Flutter
def check_flutter_exists():
  try:
    subprocess.run(['which', 'flutter'], check=True,
                   stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    return True
  except subprocess.CalledProcessError:
    return False


if check_flutter_exists():
  flutter = ("Installed")
else:
  flutter = ("Not Installed")

print("Flutter: ", flutter)

# Xcode


def check_xcode_installed():
    try:
        xcode_path = subprocess.check_output(["xcode-select", "--print-path"])
        return True
    except subprocess.CalledProcessError:
        return False


if check_xcode_installed():
    xcode = ("Installed")
else:
    xcode = ("Not Installed")

print("Xcode: ", xcode)

# VS Code


def is_vscode_installed():
    try:
        subprocess.run(["which", "code"], check=True, capture_output=True)
        return True
    except subprocess.CalledProcessError:
        return False


if is_vscode_installed():
    vscode = ("Installed")
else:
    vscode = ("Not Installed")

print("VS Code: ", vscode)


# INSTALLATION
print("")
print("")
print("INSTALLATION")
print("----------")

# if Apple Silicon Macs
if cpu == "arm":
    subprocess.call(["sudo", "softwareupdate",
                    "--install-rosetta", "--agree-to-license"])
else:
    print()

# xCLI
print()
if xCLI == "Not Installed":
    print("Installing Xcode Command Line Tools")
    subprocess.call(["xcode-select", "--install"])
else:
    print("*")

# Homebrew
print()
print()
if brew == "Not Installed":
    print("Homebrew")
    subprocess.call(
        ["/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"])
else:
    print("**")

# CocoaPods
print()
print()
if cPods == "Not Installed":
    print("Installing Cocoapods")
    subprocess.call(["brew", "install", "cocoapods"])
else:
    print("***")

# Flutter
print()
print()
if flutter == "Not Installed":
    print("Installing Flutter")
    subprocess.call(["brew", "install", "flutter"])
else:
    print("****")

# Xcode
print()
print()
if xcode == "Not Installed":
    print("Please refer to https://apps.apple.com/us/app/xcode/id497799835?mt=12 to install Xcode on your system.")
else:
    print("*****")

# VSCode
print()
print()
if vscode == "Not Installed":
    print("Installing VS Code")
    subprocess.call(["brew", "install", "--cask", "visual-studio-code"])
else:
    print("******")


# Re-Check Requirement
print()
print()
print(" ")
print("Rechecking Applications: ")
print("--------------------------")

# Xcode CLI
try:
    subprocess.run(["xcode-select", "--version"], check=True)
    xCLI = ("Installed")
except subprocess.CalledProcessError:
    xCLI = ("Not Installed")

print("Xcode Command Line Tools: ", xCLI)

# Homebrew


def is_brew_installed():
    try:
        subprocess.run(['which', 'brew'], check=True, stdout=subprocess.PIPE)
        return True
    except subprocess.CalledProcessError:
        return False


if is_brew_installed():
    brew = ("Installed.")
else:
   brew = ("Not Installed.")


print("Homebrew: ", brew)

# Cocoapods


def is_cocoapods_installed():
    try:
        subprocess.run(['which', 'pod'], check=True, stdout=subprocess.PIPE)
        return True
    except subprocess.CalledProcessError:
        return False


if is_cocoapods_installed():
    cPods = ("Installed")
else:
    cPods = ("Not Installed")

print("Cocoapods: ", cPods)


# Flutter
def check_flutter_exists():
  try:
    subprocess.run(['which', 'flutter'], check=True,
                   stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    return True
  except subprocess.CalledProcessError:
    return False


if check_flutter_exists():
  flutter = ("Installed")
else:
  flutter = ("Not Installed")

print("Flutter: ", flutter)

# Xcode


def check_xcode_installed():
    try:
        xcode_path = subprocess.check_output(["xcode-select", "--print-path"])
        return True
    except subprocess.CalledProcessError:
        return False


if check_xcode_installed():
    xcode = ("Installed")
else:
    xcode = ("Not Installed")

print("Xcode: ", xcode)

# VS Code


def is_vscode_installed():
    try:
        subprocess.run(["which", "code"], check=True, capture_output=True)
        return True
    except subprocess.CalledProcessError:
        return False


if is_vscode_installed():
    vscode = ("Installed")
else:
    vscode = ("Not Installed")

print("VS Code: ", vscode)
print("------------END------------")

print()
print()
print("~RUNNING FLUTTER DOCTOR~")
subprocess.call(["flutter", "doctor"])

print()
print()
print("IF THEY ARE ERRORS, KINDLY RE-RUN THE SCRIPT OR CONTACT US")
