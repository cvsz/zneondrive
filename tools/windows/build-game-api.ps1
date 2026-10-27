[CmdletBinding()]
param([Parameter()][string]$RepositoryRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path)
$ErrorActionPreference = "Stop"
$go = Get-Command go -ErrorAction Stop
$servicePath = Join-Path $RepositoryRoot "services\game-api"
$outputPath = Join-Path ([System.IO.Path]::GetTempPath()) ("zneondrive-api-check-" + [guid]::NewGuid().ToString("N") + ".exe")
function Invoke-Go {
    param([Parameter(Mandatory)][string[]]$Arguments)
    & $go.Source @Arguments
    if ($LASTEXITCODE -ne 0) { throw "go $($Arguments -join ' ') failed with exit code $LASTEXITCODE" }
}
Push-Location $servicePath
try {
    Invoke-Go -Arguments @("build", "-trimpath", "-o", $outputPath, "./cmd/server")
    Invoke-Go -Arguments @("vet", "./...")
    Invoke-Go -Arguments @("test", "./...")
    Write-Host "Go build, vet, and unit tests passed."
}
finally {
    Pop-Location
    if (Test-Path -LiteralPath $outputPath) { Remove-Item -LiteralPath $outputPath -Force }
}
