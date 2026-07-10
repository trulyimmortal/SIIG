# SIIG Uninstall Script (Windows PowerShell)
# Usage: irm https://get.siig.dev/uninstall.ps1 | iex

$ErrorActionPreference = "Stop"

$installDir = "$env:USERPROFILE\.siig\bin"
$binaryName = "siig.exe"
$binaryPath = Join-Path $installDir $binaryName

if (-not (Test-Path $binaryPath)) {
    Write-Host "SIIG is not installed at $binaryPath" -ForegroundColor Yellow
    Write-Host "Nothing to uninstall."
    exit 0
}

Write-Host "-> Removing SIIG..." -ForegroundColor Cyan

# Remove binary
Remove-Item $binaryPath -Force

# Remove install directory if empty
if ((Get-ChildItem $installDir -Force | Measure-Object).Count -eq 0) {
    Remove-Item $installDir -Force
}

# Remove from PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -like "*$installDir*") {
    Write-Host "-> Removing $installDir from PATH..." -ForegroundColor Cyan
    $newPath = ($userPath -split ";" | Where-Object { $_ -ne $installDir }) -join ";"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
}

Write-Host ""
Write-Host "SIIG uninstalled successfully." -ForegroundColor Green
