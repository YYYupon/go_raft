param (
    [int]$count = 3,
    [int]$basePort = 8040
)

for ($i = 0; $i -lt $count; $i++) {

    $port = $basePort + $i

    Start-Process powershell `
        -ArgumentList "-NoExit", "-Command", "./test.exe -port=$port"
}
