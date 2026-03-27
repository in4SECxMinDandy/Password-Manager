# Password Manager - Auto Start Script
# Tự động khởi động tất cả services và restart khi có lỗi

param(
    [switch]$Stop,
    [switch]$Restart
)

$ErrorActionPreference = "Continue"
$scriptPath = $MyInvocation.MyCommand.Path
if ($scriptPath) {
    $projectRoot = Split-Path -Parent $scriptPath
} else {
    $projectRoot = $PSScriptRoot
}
if (-not $projectRoot) {
    $projectRoot = Get-Location
}

# Refresh PATH to include newly installed tools
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

$backendDir = Join-Path $projectRoot "backend"
$frontendDir = Join-Path $projectRoot "frontend"

function Write-Log {
    param([string]$Message, [string]$Level = "INFO")
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $colors = @{
        "INFO" = "Cyan"
        "SUCCESS" = "Green"
        "WARN" = "Yellow"
        "ERROR" = "Red"
    }
    Write-Host "[$timestamp] [$Level] $Message" -ForegroundColor $colors[$Level]
}

function Stop-Services {
    Write-Log "Stopping all services..." "WARN"
    
    # Stop backend
    $backendProcs = Get-Process -Name "server" -ErrorAction SilentlyContinue
    if ($backendProcs) {
        $backendProcs | Stop-Process -Force
        Write-Log "Backend stopped"
    }
    
    # Stop frontend (vite/node)
    $viteProcs = Get-Process -Name "node" -ErrorAction SilentlyContinue | Where-Object { 
        $_.CommandLine -like "*vite*" -or $_.CommandLine -like "*react*"
    }
    if ($viteProcs) {
        $viteProcs | Stop-Process -Force
        Write-Log "Frontend stopped"
    }
    
    # Stop docker services
    Push-Location $projectRoot
    docker-compose down 2>$null
    Pop-Location
    
    Write-Log "All services stopped" "SUCCESS"
}

function Start-Docker-Services {
    Write-Log "Starting Docker services (PostgreSQL, Redis)..."
    
    Push-Location $projectRoot
    
    # Check if PostgreSQL is running on port 5432
    $pgPortInUse = (Get-NetTCPConnection -LocalPort 5432 -ErrorAction SilentlyContinue | Where-Object { $_.State -eq 'Listen' }).Count -gt 0
    
    if ($pgPortInUse) {
        Write-Log "PostgreSQL is already running on port 5432" "INFO"
    } else {
        # Try to start PostgreSQL container
        $pgResult = docker-compose up -d postgres 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Log "Failed to start PostgreSQL: $pgResult" "WARN"
        } else {
            Write-Log "PostgreSQL container started" "SUCCESS"
        }
    }
    
    # Check if Redis is running on port 6379
    $redisPortInUse = (Get-NetTCPConnection -LocalPort 6379 -ErrorAction SilentlyContinue | Where-Object { $_.State -eq 'Listen' }).Count -gt 0
    
    if ($redisPortInUse) {
        Write-Log "Redis is already running on port 6379" "INFO"
    } else {
        # Try to start Redis container
        $redisResult = docker-compose up -d redis 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Log "Failed to start Redis: $redisResult" "WARN"
        } else {
            Write-Log "Redis container started" "SUCCESS"
        }
    }
    
    # Wait for PostgreSQL
    Write-Log "Waiting for PostgreSQL..."
    $maxRetries = 30
    for ($i = 0; $i -lt $maxRetries; $i++) {
        $pgReady = docker exec passwordmanager_postgres pg_isready -U passwordmanager 2>$null
        if ($pgReady -match "accepting connections") {
            Write-Log "PostgreSQL is ready" "SUCCESS"
            break
        }
        Start-Sleep -Seconds 2
    }
    
    # Wait for Redis
    Write-Log "Waiting for Redis..."
    for ($i = 0; $i -lt $maxRetries; $i++) {
        $redisReady = docker exec passwordmanager_redis redis-cli ping 2>$null
        if ($redisReady -eq "PONG") {
            Write-Log "Redis is ready" "SUCCESS"
            break
        }
        Start-Sleep -Seconds 1
    }
    
    Pop-Location
    return $true
}

function Start-Backend-Service {
    param([int]$MaxRetries = 5)
    
    Write-Log "Starting Backend (Go/Fiber)..."
    
    for ($attempt = 1; $attempt -le $MaxRetries; $attempt++) {
        # Start backend in background
        $processInfo = New-Object System.Diagnostics.ProcessStartInfo
        $processInfo.FileName = "go"
        $processInfo.Arguments = "run ./cmd/server"
        $processInfo.WorkingDirectory = $backendDir
        $processInfo.UseShellExecute = $false
        $processInfo.RedirectStandardOutput = $true
        $processInfo.RedirectStandardError = $true
        $processInfo.CreateNoWindow = $true
        
        $process = New-Object System.Diagnostics.Process
        $process.StartInfo = $processInfo
        $process.Start() | Out-Null
        
        Start-Sleep -Seconds 5
        
        # Check if backend is running
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -TimeoutSec 5 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-Log "Backend is running on http://localhost:8080" "SUCCESS"
                return $process
            }
        } catch {
            # Not ready yet, check stderr
            $stderr = $process.StandardError.ReadToEnd()
            if ($stderr -and $attempt -eq $MaxRetries) {
                Write-Log "Backend error: $stderr" "ERROR"
            }
        }
        
        # Clean up failed attempt
        if (-not $process.HasExited) {
            $process.Kill()
        }
        $process.Dispose()
        
        Write-Log "Backend attempt $attempt failed, retrying..." "WARN"
        Start-Sleep -Seconds 3
    }
    
    return $null
}

function Start-Frontend-Service {
    param([int]$MaxRetries = 5)
    
    Write-Log "Starting Frontend (React/Vite)..."
    
    for ($attempt = 1; $attempt -le $MaxRetries; $attempt++) {
        # Check if dependencies are installed
        $nodeModulesPath = Join-Path $frontendDir "node_modules"
        if (-not (Test-Path $nodeModulesPath)) {
            Write-Log "Installing frontend dependencies..."
            $npmStartInfo = New-Object System.Diagnostics.ProcessStartInfo
            $npmStartInfo.FileName = "npm"
            $npmStartInfo.Arguments = "install"
            $npmStartInfo.WorkingDirectory = $frontendDir
            $npmStartInfo.UseShellExecute = $false
            $npmStartInfo.RedirectStandardOutput = $true
            $npmStartInfo.RedirectStandardError = $true
            $npmStartInfo.CreateNoWindow = $true
            
            $npmProcess = New-Object System.Diagnostics.Process
            $npmProcess.StartInfo = $npmStartInfo
            $npmProcess.Start() | Out-Null
            $npmProcess.WaitForExit()
            
            if ($npmProcess.ExitCode -ne 0) {
                Write-Log "Failed to install dependencies" "ERROR"
                $npmProcess.Dispose()
                return $null
            }
            $npmProcess.Dispose()
        }
        
        # Start frontend
        $viteStartInfo = New-Object System.Diagnostics.ProcessStartInfo
        $viteStartInfo.FileName = "npm"
        $viteStartInfo.Arguments = "run dev"
        $viteStartInfo.WorkingDirectory = $frontendDir
        $viteStartInfo.UseShellExecute = $false
        $viteStartInfo.RedirectStandardOutput = $true
        $viteStartInfo.RedirectStandardError = $true
        $viteStartInfo.CreateNoWindow = $true
        
        $viteProcess = New-Object System.Diagnostics.Process
        $viteProcess.StartInfo = $viteStartInfo
        $viteProcess.Start() | Out-Null
        
        Start-Sleep -Seconds 5
        
        # Check if frontend is running
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:5173" -TimeoutSec 5 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-Log "Frontend is running on http://localhost:5173" "SUCCESS"
                return $viteProcess
            }
        } catch {
            # Not ready yet
        }
        
        # Clean up failed attempt
        if (-not $viteProcess.HasExited) {
            $viteProcess.Kill()
        }
        $viteProcess.Dispose()
        
        Write-Log "Frontend attempt $attempt failed, retrying..." "WARN"
        Start-Sleep -Seconds 3
    }
    
    return $null
}

function Watch-And-Restart {
    param($backendJob, $frontendJob)
    
    Write-Log "Starting watchdog to monitor services..." "INFO"
    Write-Log "Press Ctrl+C to stop all services" "INFO"
    
    while ($true) {
        Start-Sleep -Seconds 10
        
        # Check backend
        $backendRunning = $false
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -TimeoutSec 3 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                $backendRunning = $true
            }
        } catch {
            $backendRunning = $false
        }
        
        # Check frontend
        $frontendRunning = $false
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:5173" -TimeoutSec 3 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                $frontendRunning = $true
            }
        } catch {
            $frontendRunning = $false
        }
        
        # Restart backend if down
        if (-not $backendRunning) {
            Write-Log "Backend is down, restarting..." "WARN"
            
            # Stop old job
            if ($backendJob) {
                Stop-Job -Id $backendJob.Id -ErrorAction SilentlyContinue
                Remove-Job -Id $backendJob.Id -ErrorAction SilentlyContinue
            }
            
            $backendJob = Start-Backend-Service
            if (-not $backendJob) {
                Write-Log "Failed to restart backend after multiple attempts" "ERROR"
            }
        }
        
        # Restart frontend if down
        if (-not $frontendRunning) {
            Write-Log "Frontend is down, restarting..." "WARN"
            
            # Stop old job
            if ($frontendJob) {
                Stop-Job -Id $frontendJob.Id -ErrorAction SilentlyContinue
                Remove-Job -Id $frontendJob.Id -ErrorAction SilentlyContinue
            }
            
            $frontendJob = Start-Frontend-Service
            if (-not $frontendJob) {
                Write-Log "Failed to restart frontend after multiple attempts" "ERROR"
            }
        }
    }
}

# Main execution
Write-Host ""
Write-Host "========================================" -ForegroundColor Magenta
Write-Host "   Password Manager - Auto Start" -ForegroundColor Magenta
Write-Host "========================================" -ForegroundColor Magenta
Write-Host ""

if ($Stop) {
    Stop-Services
    exit 0
}

if ($Restart) {
    Stop-Services
    Start-Sleep -Seconds 2
}

# Stop existing services first
Stop-Services
Start-Sleep -Seconds 2

# Start Docker services
if (-not (Start-Docker-Services)) {
    Write-Log "Failed to start Docker services. Please check if Docker is running." "ERROR"
    exit 1
}

# Start backend
$backendJob = Start-Backend-Service
if (-not $backendJob) {
    Write-Log "Failed to start backend after multiple attempts" "ERROR"
    Write-Log "Check Go installation and backend code" "ERROR"
}

# Start frontend
$frontendJob = Start-Frontend-Service
if (-not $frontendJob) {
    Write-Log "Failed to start frontend after multiple attempts" "ERROR"
    Write-Log "Check Node.js installation and frontend code" "ERROR"
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "   All Services Started!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host "Backend:  http://localhost:8080" -ForegroundColor Cyan
Write-Host "Frontend: http://localhost:5173" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

# Start watchdog
try {
    Watch-And-Restart -backendJob $backendJob -frontendJob $frontendJob
} finally {
    Write-Log "Shutting down..." "WARN"
    Stop-Services
}
