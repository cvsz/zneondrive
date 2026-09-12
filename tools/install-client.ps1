param(
  [ValidateSet("package","source","play","doctor")]
  [string]$Mode = "package",
  [string]$Package = "",
  [string]$InstallDir = "$env:LOCALAPPDATA\ZeaZDev\NeonDrive",
  [string]$ApiUrl = "http://127.0.0.1:18080",
  [string]$UERoot = $env:UE_ROOT,
  [switch]$CreateShortcut
)

$ErrorActionPreference = "Stop"
$ProjectRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$UProject = Join-Path $ProjectRoot "game\NeonDrive.uproject"

function Assert-ApiUrl {
  if ($ApiUrl -notmatch '^https?://[^\s]+$') { throw "ApiUrl must be an http(s) URL." }
}

function Find-ClientExe {
  if (-not (Test-Path $InstallDir)) { return $null }
  Get-ChildItem -Path $InstallDir -Filter "NeonDrive.exe" -File -Recurse -ErrorAction SilentlyContinue |
    Select-Object -First 1
}

function Install-PlayerPackage {
  if ([string]::IsNullOrWhiteSpace($Package)) { throw "Package is required in package mode." }
  if (-not (Test-Path $Package)) { throw "Package not found: $Package" }

  if (Test-Path $InstallDir) { Remove-Item -Recurse -Force $InstallDir }
  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
  $item = Get-Item $Package

  if ($item.PSIsContainer) {
    Copy-Item -Path (Join-Path $item.FullName "*") -Destination $InstallDir -Recurse -Force
  } elseif ($item.Extension -ieq ".zip") {
    Expand-Archive -Path $item.FullName -DestinationPath $InstallDir -Force
  } else {
    throw "Windows installer accepts a packaged client directory or .zip file."
  }

  $exe = Find-ClientExe
  if (-not $exe) { throw "Package installed but NeonDrive.exe was not found." }

  if ($CreateShortcut) {
    $desktop = [Environment]::GetFolderPath("Desktop")
    $shortcutPath = Join-Path $desktop "PROJECT NEON DRIVE.lnk"
    $shell = New-Object -ComObject WScript.Shell
    $shortcut = $shell.CreateShortcut($shortcutPath)
    $shortcut.TargetPath = $exe.FullName
    $shortcut.Arguments = "-ZNeonApi=$ApiUrl"
    $shortcut.WorkingDirectory = $exe.DirectoryName
    $shortcut.Save()
    Write-Host "Shortcut created: $shortcutPath"
  }

  Write-Host "Player client installed: $($exe.FullName)"
  Write-Host "Run with -Mode play to launch it."
}

function Build-SourceClient {
  if ([string]::IsNullOrWhiteSpace($UERoot)) { throw "UERoot or UE_ROOT is required for source mode." }
  $generate = Join-Path $UERoot "Engine\Build\BatchFiles\GenerateProjectFiles.bat"
  $build = Join-Path $UERoot "Engine\Build\BatchFiles\Build.bat"
  if (-not (Test-Path $generate)) { throw "GenerateProjectFiles.bat not found under UERoot." }
  if (-not (Test-Path $build)) { throw "Build.bat not found under UERoot." }

  & $generate "-project=$UProject" -game
  if ($LASTEXITCODE -ne 0) { throw "GenerateProjectFiles failed." }
  & $build NeonDriveClient Win64 Development $UProject -WaitMutex
  if ($LASTEXITCODE -ne 0) { throw "NeonDriveClient build failed." }

  Write-Host "NeonDriveClient source target built."
  Write-Host "This is source-build readiness only; packaged-play evidence remains a separate release gate."
}

function Play-Client {
  Assert-ApiUrl
  $exe = Find-ClientExe
  if (-not $exe) { throw "No installed NeonDrive.exe found under $InstallDir." }
  $env:ZNEON_GAME_API_URL = $ApiUrl
  & $exe.FullName "-ZNeonApi=$ApiUrl"
}

function Doctor {
  Write-Host "PROJECT: NEON DRIVE Windows client doctor"
  Write-Host "Project root: $ProjectRoot"
  Write-Host "Install dir:  $InstallDir"
  Write-Host "API URL:      $ApiUrl"
  Write-Host "PowerShell:   $($PSVersionTable.PSVersion)"
  if ($UERoot) { Write-Host "UE_ROOT:      $UERoot" } else { Write-Host "UE_ROOT:      not set" }
  $exe = Find-ClientExe
  if ($exe) { Write-Host "Client:       $($exe.FullName)" } else { Write-Host "Client:       not installed" }
  Write-Host "Server shared keys are intentionally not accepted by this client installer."
}

Assert-ApiUrl
switch ($Mode) {
  "package" { Install-PlayerPackage }
  "source"  { Build-SourceClient }
  "play"    { Play-Client }
  "doctor"  { Doctor }
}
