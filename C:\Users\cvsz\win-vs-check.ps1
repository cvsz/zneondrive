Get-ChildItem 'HKLM:SOFTWARE\Microsoft\VisualStudio' -ErrorAction SilentlyContinue | Get-ItemProperty -ErrorAction SilentlyContinue | Select-Object InstallationPath, InstallVersion, ProductId | Format-Table -AutoSize
Get-ChildItem 'C:\Program Files\Microsoft Visual Studio' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Name
Get-ChildItem 'C:\Program Files (x86)\Microsoft Visual Studio' -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Name
Get-Command msbuild.exe 2>&1
Get-Command cl.exe 2>&1
