# Extract UE zip on Windows
$startTime = Get-Date
Write-Host "Extracting UE 5.8.2 zip..."
$zipPath = "C:\Users\cvsz\Linux_Unreal_Engine_5.8.2.zip"
$destPath = "D:\UnrealEngine-5.8"
Remove-Item $destPath -Recurse -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Path $destPath -Force | Out-Null

# Use 7z if available, otherwise Expand-Archive
if (Get-Command 7z -ErrorAction SilentlyContinue) {
    Write-Host "Using 7z for extraction..."
    7z x $zipPath -o$destPath -y 2>&1
} else {
    Write-Host "Using Expand-Archive..."
    Expand-Archive -Path $zipPath -DestinationPath $destPath -Force
}

Write-Host "Extraction complete at $(Get-Date), took $(($(Get-Date) - $startTime))"
