$url = "http://192.168.1.85:8765/Linux_Unreal_Engine_5.8.2.zip"
$out = "D:\Data\Linux_Unreal_Engine_5.8.2.zip"
Write-Host "Starting download from $url"
$startTime = Get-Date
Invoke-WebRequest -Uri $url -OutFile $out -UseBasicParsing -Verbose 2>&1
Write-Host "Download complete at $(Get-Date), took $(($(Get-Date) - $startTime))"
