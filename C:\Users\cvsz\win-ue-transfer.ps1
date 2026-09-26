# Transfer and extract UE 5.8
Write-Host "Starting UE transfer..."
$startTime = Get-Date
# Remove existing zip if present
Remove-Item D:\Data\Linux_Unreal_Engine_5.8.2.zip -ErrorAction SilentlyContinue
# Copy from Linux via SMB/network path or use the existing one
# First check if we can access it
if (Test-Path "\\192.168.1.85\mnt\zworkforce-storage\Linux_Unreal_Engine_5.8.2.zip") {
    Write-Host "Copying from Linux SMB share..."
    Copy-Item "\\192.168.1.85\mnt\zworkforce-storage\Linux_Unreal_Engine_5.8.2.zip" D:\Data\
} else {
    Write-Host "SMB share not available, will use scp/rsync"
}
Write-Host "Transfer time: $((Get-Date) - $startTime)"
