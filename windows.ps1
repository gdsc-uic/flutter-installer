"##########################################################"
"#        GDSC-UIC Resoruce Installer for Flutter         #" 
"##########################################################"
<#check and install chocolatey#>
try {
    choco | Out-Null
    "Chocolatey: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
    "Chocolatey: Installed"
}

refreshenv


<#check and install Git#>
try {
    git | Out-Null
    "Git: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install git
    "Git: Installed"
}

<#check and install Dart-SDK#>
try {
    dart | Out-Null
    "Dart-SDK: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install dart-sdk
    "Dart-SDK: Installed"
}

<#check and install JAVASDK#>
try {
    javac | Out-Null
    "JDK: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install jdk11
    "Java-SDK: Installed"
}

<#check and install flutter#>
try {
    flutter | Out-Null
    "Flutter: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install flutter
    "Flutter: Installed"
}

<#check and install Android Studio#>
function check_if_installed($p1) {
    $p1 = "Android Studio"
    $software = $p1;
    $installed = ((gp HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*).DisplayName -Match $p1).Length -gt 0
    If (-Not $installed) {
        Write-Host "Installing '$software'.";
        choco install androidstudio -y -f
        refreshenv
        Write-Host "'$software'  is not installed.";
    }
    else {
        Write-Host "'$software' is installed."
    }
    return 
}
check_if_installed -p1 "Visual Studio Code"

<#check and install VSCode#>
function check_if_installed($p1) {
    $p1 = "Visual Studio Code"
    $software = $p1;
    $installed = ((gp HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*).DisplayName -Match $p1).Length -gt 0
    If (-Not $installed) {
        Write-Host "Installing '$software'."
        choco install vscode -y -f
        refreshenv
        Write-Host "'$software'  is not installed.";
    }
    else {
        Write-Host "'$software' is installed."
    }
    return
}
check_if_installed -p1 "Visual Studio Code"

<#Checker#>
" "
" "
"##########################################################"
"#        GDSC-UIC Resoruce Checker for Flutter           #" 
"##########################################################"
<#check and install chocolatey#>
try {
    choco | Out-Null
    "Chocolatey: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    "Chocolatey: Missing"
}


<#check and install Git#>
try {
    git | Out-Null
    "Git: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install git
    "Git: Missing"
}

<#check and install Dart-SDK#>
try {
    dart | Out-Null
    "Dart-SDK: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install dart-sdk
    "Dart-SDK: Missing"
}

<#check and install JAVASDK#>
try {
    javac | Out-Null
    "JDK: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install jdk11
    "Java-SDK: Missing"
}

<#check and install flutter#>
try {
    flutter | Out-Null
    "Flutter: Present"
}
catch [System.Management.Automation.CommandNotFoundException] {
    choco install flutter
    "Flutter: Missing"
}

<#check and install Android Studio#>
function check_if_installed($p1) {
    $p1 = "Android Studio"
    $software = $p1;
    $installed = ((gp HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*).DisplayName -Match $p1).Length -gt 0
    If (-Not $installed) {
        Write-Host "'$software'  is missing.";
    }
    else {
        Write-Host "'$software' is installed."
    }
    return 
}
check_if_installed -p1 "Visual Studio Code"

<#check and install VSCode#>
function check_if_installed($p1) {
    $p1 = "Visual Studio Code"
    $software = $p1;
    $installed = ((gp HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*).DisplayName -Match $p1).Length -gt 0
    If (-Not $installed) {
        Write-Host "'$software'  is missing.";
    }
    else {
        Write-Host "'$software' is installed."
    }
    return
}
check_if_installed -p1 "Visual Studio Code"
" "
" If some reported as 'Missing' or 'Not Installed', Kindly re-run the script or contact us."
" "
" If all statements are 'Present and Installed' then close the Script and run 'flutter doctor -v'"
" " 