Get-ChildItem "$PSScriptRoot\*.json" | ForEach-Object {
  try {
    Get-Content $_.FullName -Raw | ConvertFrom-Json | Out-Null
    Write-Output "OK: $($_.Name)"
  } catch {
    Write-Output "FAIL: $($_.Name)"
  }
}
