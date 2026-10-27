[CmdletBinding(SupportsShouldProcess)]
param(
    [Parameter(Mandatory)][uri]$Uri,
    [Parameter(Mandatory)][string]$DestinationPath,
    [Parameter(Mandatory)][ValidatePattern("^[A-Fa-f0-9]{64}$")][string]$Sha256,
    [Parameter()][switch]$AllowInsecureHttp
)
$ErrorActionPreference = "Stop"
if ($Uri.Scheme -ne "https" -and -not ($AllowInsecureHttp -and $Uri.Scheme -eq "http")) {
    throw "HTTPS is required. Use -AllowInsecureHttp only for a trusted isolated transfer with an independently verified SHA-256."
}
$destination = [System.IO.Path]::GetFullPath($DestinationPath)
$parent = Split-Path -Parent $destination
if (-not (Test-Path -LiteralPath $parent -PathType Container)) { New-Item -ItemType Directory -Path $parent -Force | Out-Null }
if (Test-Path -LiteralPath $destination) { throw "Destination already exists; choose a new path to avoid overwriting an existing artifact." }
if ($PSCmdlet.ShouldProcess($destination, "Download and verify SHA-256 from $Uri")) {
    $temporary = Join-Path $parent (".ue-download-" + [guid]::NewGuid().ToString("N") + ".tmp")
    try {
        Invoke-WebRequest -Uri $Uri -OutFile $temporary
        $actual = (Get-FileHash -LiteralPath $temporary -Algorithm SHA256).Hash
        if ($actual -ne $Sha256.ToUpperInvariant()) { throw "SHA-256 mismatch. Expected $($Sha256.ToUpperInvariant()); received $actual." }
        Move-Item -LiteralPath $temporary -Destination $destination
        Write-Host "Verified artifact: $destination"
        Write-Host "SHA256: $actual"
    }
    finally {
        if (Test-Path -LiteralPath $temporary) { Remove-Item -LiteralPath $temporary -Force }
    }
}
