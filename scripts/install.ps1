# SIIG Install Script (Windows PowerShell)
# Usage: irm https://get.siig.dev/install.ps1 | iex

$ErrorActionPreference = "Stop"

$repo = "yv3000/SIIG"
$binaryName = "siig.exe"
$installDir = "$env:USERPROFILE\.siig\bin"

Write-Host "-> Detecting system..." -ForegroundColor Cyan

$arch = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture) {
    "Arm64" { "arm64" }
    default { "amd64" }
}
$assetName = "siig-windows-${arch}.exe"

# Get latest release
Write-Host "-> Fetching latest release..." -ForegroundColor Cyan
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
$version = $release.tag_name

Write-Host "-> Latest version: $version" -ForegroundColor Cyan

$downloadUrl = "https://github.com/$repo/releases/download/$version/$assetName"
$checksumsUrl = "https://github.com/$repo/releases/download/$version/checksums.txt"

# Create install directory
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

$binaryPath = Join-Path $installDir $binaryName
$tempPath = Join-Path $env:TEMP $assetName

# Download binary
Write-Host "-> Downloading $assetName..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempPath -UseBasicParsing

# Verify checksum
Write-Host "-> Verifying checksum..." -ForegroundColor Cyan
try {
    $checksums = Invoke-WebRequest -Uri $checksumsUrl -UseBasicParsing
    $expectedLine = ($checksums.Content -split "`n") | Where-Object { $_ -match $assetName }
    if ($expectedLine) {
        $expected = ($expectedLine -split "\s+")[0]
        $actual = (Get-FileHash -Path $tempPath -Algorithm SHA256).Hash.ToLower()
        if ($expected -ne $actual) {
            Write-Host "Error: Checksum mismatch!" -ForegroundColor Red
            Write-Host "  Expected: $expected"
            Write-Host "  Actual:   $actual"
            Remove-Item $tempPath -Force
            exit 1
        }
        Write-Host "-> Checksum verified" -ForegroundColor Green
    } else {
        Write-Host "Warning: Could not find checksum for this asset" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Warning: Could not verify checksum" -ForegroundColor Yellow
}

# Move binary to install directory
Move-Item -Path $tempPath -Destination $binaryPath -Force

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host "-> Adding $installDir to PATH..." -ForegroundColor Cyan
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    $env:Path = "$env:Path;$installDir"
}

Write-Host ""
Write-Host "SIIG $version installed successfully!" -ForegroundColor Green
Write-Host "  Location: $binaryPath"
Write-Host ""
Write-Host "  Run 'siig' in any project directory to get started."
Write-Host ""
Write-Host "  Note: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
