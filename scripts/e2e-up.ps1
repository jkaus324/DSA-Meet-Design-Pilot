# scripts/e2e-up.ps1
# Bootstrap the dashboard for e2e runs without touching the user's progress.json.
#
# What it does:
#   1. Creates an empty e2e/.tmp-progress.json
#   2. Starts the dashboard with PROGRESS_JSON_PATH pointed at that file
#   3. Waits until /api/runner-status responds (max 10 s)
#   4. Prints the runner status and the PID so the caller can stop it later
#
# Usage:
#   pwsh scripts/e2e-up.ps1            # default: blocks in foreground
#   pwsh scripts/e2e-up.ps1 -Detach    # spawn server detached, return PID
#
# Stop the server with:  Stop-Process -Id <PID>

param(
    [switch]$Detach
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$progressFile = Join-Path $repoRoot 'e2e/.tmp-progress.json'
$serverScript = Join-Path $repoRoot 'dashboard/server.js'

# 1. Ensure the temp progress file exists (clean state every run).
$progressDir = Split-Path -Parent $progressFile
if (-not (Test-Path $progressDir)) { New-Item -ItemType Directory -Path $progressDir | Out-Null }
'{"version":3,"problems":{}}' | Set-Content -Encoding utf8 -Path $progressFile
Write-Host "[e2e-up] Reset $progressFile"

# 2. Start the server with the env override.
$env:PROGRESS_JSON_PATH = 'e2e/.tmp-progress.json'

if ($Detach) {
    $proc = Start-Process -FilePath 'node' -ArgumentList $serverScript `
        -WorkingDirectory $repoRoot -PassThru -WindowStyle Hidden
    Write-Host "[e2e-up] Server started (PID $($proc.Id)) — detached"
} else {
    Write-Host "[e2e-up] Starting server in foreground (Ctrl+C to stop)..."
    Push-Location $repoRoot
    try { node $serverScript } finally { Pop-Location }
    return
}

# 3. Poll /api/runner-status until ready.
$deadline = (Get-Date).AddSeconds(10)
$ready = $false
while ((Get-Date) -lt $deadline) {
    try {
        $resp = Invoke-RestMethod -Uri 'http://localhost:3000/api/runner-status' -TimeoutSec 1
        $ready = $true
        break
    } catch { Start-Sleep -Milliseconds 250 }
}

if (-not $ready) {
    Write-Error "[e2e-up] Server did not respond on :3000 within 10s. Killing PID $($proc.Id)."
    Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
    exit 1
}

# 4. Report toolchain availability so callers can decide whether to skip.
Write-Host "[e2e-up] Runner status:"
Write-Host "         cpp_available = $($resp.cpp_available)"
Write-Host "         go_available  = $($resp.go_available)"
Write-Host "[e2e-up] PID  = $($proc.Id)"
Write-Host "[e2e-up] Stop with:  Stop-Process -Id $($proc.Id)"
