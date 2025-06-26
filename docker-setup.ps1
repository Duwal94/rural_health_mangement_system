# Rural Health Management System - Docker Setup Script (PowerShell)

param(
    [Parameter(Position=0)]
    [ValidateSet("start", "seed", "status", "logs", "stop", "restart", "clean", "help")]
    [string]$Command = "start"
)

function Write-Header {
    Write-Host "🏥 Rural Health Management System - Docker Setup" -ForegroundColor Cyan
    Write-Host "================================================" -ForegroundColor Cyan
}

function Test-Docker {
    Write-Host "Checking Docker installation..." -ForegroundColor Yellow
    
    if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Host "❌ Docker is not installed. Please install Docker Desktop first." -ForegroundColor Red
        Write-Host "   Visit: https://docs.docker.com/desktop/windows/" -ForegroundColor Yellow
        exit 1
    }
    
    if (!(Get-Command docker-compose -ErrorAction SilentlyContinue)) {
        Write-Host "❌ Docker Compose is not installed. Please install Docker Compose first." -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✅ Docker and Docker Compose are installed" -ForegroundColor Green
}

function Build-App {
    Write-Host "🔨 Building the application..." -ForegroundColor Yellow
    docker-compose build --no-cache
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Application built successfully" -ForegroundColor Green
    } else {
        Write-Host "❌ Build failed" -ForegroundColor Red
        exit 1
    }
}

function Start-Services {
    Write-Host "🚀 Starting services..." -ForegroundColor Yellow
    docker-compose up -d
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Services started successfully" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to start services" -ForegroundColor Red
        exit 1
    }
}

function Test-Health {
    Write-Host "🔍 Checking service health..." -ForegroundColor Yellow
    
    # Wait for database
    Write-Host "   Waiting for database..." -ForegroundColor Yellow
    for ($i = 1; $i -le 30; $i++) {
        try {
            docker-compose exec -T postgres pg_isready -U postgres 2>$null
            if ($LASTEXITCODE -eq 0) {
                Write-Host "   ✅ Database is ready" -ForegroundColor Green
                break
            }
        } catch {}
        
        Write-Host "   ⏳ Waiting for database... ($i/30)" -ForegroundColor Yellow
        Start-Sleep -Seconds 2
    }
    
    # Wait for API
    Write-Host "   Waiting for API..." -ForegroundColor Yellow
    for ($i = 1; $i -le 30; $i++) {
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:3000/health" -UseBasicParsing -TimeoutSec 2 2>$null
            if ($response.StatusCode -eq 200) {
                Write-Host "   ✅ API is ready" -ForegroundColor Green
                break
            }
        } catch {}
        
        Write-Host "   ⏳ Waiting for API... ($i/30)" -ForegroundColor Yellow
        Start-Sleep -Seconds 2
    }
}

function Add-SeedData {
    Write-Host "🌱 Seeding database with sample data..." -ForegroundColor Yellow
    docker-compose exec api go run cmd/seed/main.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Database seeded successfully" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to seed database" -ForegroundColor Red
    }
}

function Show-Status {
    Write-Host "📊 Service Status:" -ForegroundColor Cyan
    docker-compose ps
    Write-Host ""
    Write-Host "🌐 Application URLs:" -ForegroundColor Cyan
    Write-Host "   API: http://localhost:3000" -ForegroundColor White
    Write-Host "   Health Check: http://localhost:3000/health" -ForegroundColor White
    Write-Host "   Database: localhost:5432" -ForegroundColor White
}

function Show-Logs {
    Write-Host "📝 Recent logs:" -ForegroundColor Cyan
    docker-compose logs --tail=20
}

function Stop-Services {
    Write-Host "🛑 Stopping services..." -ForegroundColor Yellow
    docker-compose down
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Services stopped" -ForegroundColor Green
    }
}

function Remove-Everything {
    Write-Host "🧹 Cleaning up..." -ForegroundColor Yellow
    docker-compose down -v --remove-orphans
    docker system prune -f
    Write-Host "✅ Cleanup completed" -ForegroundColor Green
}

function Show-Help {
    Write-Host "Usage: .\docker-setup.ps1 [command]" -ForegroundColor White
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Cyan
    Write-Host "  start    - Build and start all services (default)" -ForegroundColor White
    Write-Host "  seed     - Seed database with sample data" -ForegroundColor White
    Write-Host "  status   - Show service status" -ForegroundColor White
    Write-Host "  logs     - Show recent logs" -ForegroundColor White
    Write-Host "  stop     - Stop all services" -ForegroundColor White
    Write-Host "  restart  - Restart all services" -ForegroundColor White
    Write-Host "  clean    - Stop services and clean up" -ForegroundColor White
    Write-Host "  help     - Show this help message" -ForegroundColor White
}

# Main script execution
Write-Header

switch ($Command) {
    "start" {
        Test-Docker
        Build-App
        Start-Services
        Test-Health
        Show-Status
        Write-Host ""
        Write-Host "🎉 Rural Health Management System is now running!" -ForegroundColor Green
        Write-Host "   Try: Invoke-WebRequest http://localhost:3000/health" -ForegroundColor Yellow
    }
    "seed" {
        Add-SeedData
    }
    "status" {
        Show-Status
    }
    "logs" {
        Show-Logs
    }
    "stop" {
        Stop-Services
    }
    "restart" {
        Stop-Services
        Start-Sleep -Seconds 2
        Start-Services
        Test-Health
        Show-Status
    }
    "clean" {
        Remove-Everything
    }
    "help" {
        Show-Help
    }
}
