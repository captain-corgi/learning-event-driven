# Setup script for Learning Event-Driven Microservices (PowerShell)
# This script helps initialize the development environment on Windows

Write-Host "🚀 Setting up Learning Event-Driven Microservices environment..." -ForegroundColor Green

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "✅ Go detected: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go is not installed. Please install Go 1.21+ first." -ForegroundColor Red
    Write-Host "   Visit: https://golang.org/doc/install" -ForegroundColor Yellow
    exit 1
}

# Initialize Go module if not exists
if (-not (Test-Path "go.mod")) {
    Write-Host "📦 Initializing Go module..." -ForegroundColor Blue
    go mod init learning-event-driven
}

# Download dependencies
Write-Host "📥 Downloading Go dependencies..." -ForegroundColor Blue
go mod tidy

# Check if Docker is installed
try {
    docker --version | Out-Null
    Write-Host "✅ Docker detected" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Docker not found. You'll need Docker for containerization modules." -ForegroundColor Yellow
    Write-Host "   Visit: https://docs.docker.com/get-docker/" -ForegroundColor Yellow
}

# Check if kubectl is installed
try {
    kubectl version --client | Out-Null
    Write-Host "✅ kubectl detected" -ForegroundColor Green
} catch {
    Write-Host "⚠️  kubectl not found. You'll need kubectl for Kubernetes modules." -ForegroundColor Yellow
    Write-Host "   Visit: https://kubernetes.io/docs/tasks/tools/" -ForegroundColor Yellow
}

# Create module directories if they don't exist
Write-Host "📁 Creating module directories..." -ForegroundColor Blue
for ($i = 1; $i -le 10; $i++) {
    $moduleDir = "modules/module-{0:D2}" -f $i
    if (-not (Test-Path $moduleDir)) {
        New-Item -ItemType Directory -Path $moduleDir -Force | Out-Null
        Write-Host "   Created $moduleDir" -ForegroundColor Gray
    }
}

# Create shared directories
Write-Host "📁 Creating shared directories..." -ForegroundColor Blue
$sharedDirs = @(
    "shared/utils",
    "shared/types", 
    "shared/config",
    "deployments/kubernetes",
    "deployments/helm",
    "deployments/docker",
    "scripts/build",
    "scripts/deploy",
    "scripts/test"
)

foreach ($dir in $sharedDirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
    }
}

# Create .env.example file
Write-Host "📝 Creating environment template..." -ForegroundColor Blue
$envContent = @"
# Environment Configuration Template
# Copy this file to .env and update with your values

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=eventdriven
DB_USER=postgres
DB_PASSWORD=password

# Message Broker
KAFKA_BROKERS=localhost:9092
NATS_URL=nats://localhost:4222

# Observability
JAEGER_ENDPOINT=http://localhost:14268/api/traces
PROMETHEUS_URL=http://localhost:9090

# Application
LOG_LEVEL=info
HTTP_PORT=8080
GRPC_PORT=9090
"@

$envContent | Out-File -FilePath ".env.example" -Encoding UTF8

Write-Host ""
Write-Host "🎉 Setup complete! Next steps:" -ForegroundColor Green
Write-Host ""
Write-Host "1. Review the learning path: docs/learning-path.md" -ForegroundColor White
Write-Host "2. Check Git workflow: docs/git-flow-strategy.md" -ForegroundColor White
Write-Host "3. Create develop branch:" -ForegroundColor White
Write-Host "   git checkout -b develop" -ForegroundColor Gray
Write-Host "   git push -u origin develop" -ForegroundColor Gray
Write-Host ""
Write-Host "4. Start with Module 1:" -ForegroundColor White
Write-Host "   git checkout -b feature/module-1-foundations" -ForegroundColor Gray
Write-Host "   cd modules/module-01" -ForegroundColor Gray
Write-Host ""
Write-Host "Happy learning! 🚀" -ForegroundColor Green
