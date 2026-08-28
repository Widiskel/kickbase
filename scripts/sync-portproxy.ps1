# Sync all Podman running container ports to Windows localhost via netsh portproxy
$wslIp = (wsl -d podman-machine-default ip addr show eth0 | Select-String -Pattern 'inet (\d+\.\d+\.\d+\.\d+)' | ForEach-Object { $_.Matches.Groups[1].Value })

if (-not $wslIp) {
    Write-Error "Could not determine Podman WSL IP"
    exit 1
}

Write-Host "Syncing Podman ports to Localhost (WSL IP: $wslIp)..."

# Get all published ports from podman ps
$portsOutput = podman ps --format "{{.Ports}}"
$ports = @()

foreach ($line in $portsOutput) {
    # Match patterns like 0.0.0.0:8080->8080/tcp or :::8080->8080/tcp
    $matches = [regex]::Matches($line, ':(\d+)->(\d+)')
    foreach ($m in $matches) {
        $hostPort = $m.Groups[1].Value
        $containerPort = $m.Groups[2].Value
        if ($hostPort -notin $ports) {
            $ports += $hostPort
        }
    }
}

foreach ($p in $ports) {
    Write-Host "Forwarding localhost:$p -> $wslIp`:$p"
    Start-Process netsh -ArgumentList "interface portproxy add v4tov4 listenport=$p listenaddress=127.0.0.1 connectport=$p connectaddress=$wslIp" -Verb RunAs -WindowStyle Hidden
}

Write-Host "Successfully synced all container ports ($($ports -join ', ')) to localhost!"
