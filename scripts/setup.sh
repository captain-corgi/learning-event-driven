#!/bin/bash

# Setup script for Learning Event-Driven Microservices
# This script helps initialize the development environment

set -e

echo "🚀 Setting up Learning Event-Driven Microservices environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please upgrade to Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Initialize Go module if not exists
if [ ! -f "go.mod" ]; then
    echo "📦 Initializing Go module..."
    go mod init learning-event-driven
fi

# Download dependencies
echo "📥 Downloading Go dependencies..."
go mod tidy

# Check if Docker is installed
if command -v docker &> /dev/null; then
    echo "✅ Docker detected"
else
    echo "⚠️  Docker not found. You'll need Docker for containerization modules."
    echo "   Visit: https://docs.docker.com/get-docker/"
fi

# Check if kubectl is installed
if command -v kubectl &> /dev/null; then
    echo "✅ kubectl detected"
else
    echo "⚠️  kubectl not found. You'll need kubectl for Kubernetes modules."
    echo "   Visit: https://kubernetes.io/docs/tasks/tools/"
fi

# Create module directories if they don't exist
echo "📁 Creating module directories..."
for i in {01..10}; do
    module_dir="modules/module-$i"
    if [ ! -d "$module_dir" ]; then
        mkdir -p "$module_dir"
        echo "   Created $module_dir"
    fi
done

# Create shared directories
echo "📁 Creating shared directories..."
mkdir -p shared/{utils,types,config}
mkdir -p deployments/{kubernetes,helm,docker}
mkdir -p scripts/{build,deploy,test}

# Set up Git hooks (optional)
if [ -d ".git" ]; then
    echo "🔧 Setting up Git hooks..."
    
    # Pre-commit hook
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Pre-commit hook for Go projects

echo "Running pre-commit checks..."

# Format code
go fmt ./...

# Run tests
if ! go test ./...; then
    echo "❌ Tests failed. Commit aborted."
    exit 1
fi

# Run linter if available
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run; then
        echo "❌ Linting failed. Commit aborted."
        exit 1
    fi
fi

echo "✅ Pre-commit checks passed"
EOF
    
    chmod +x .git/hooks/pre-commit
    echo "   Pre-commit hook installed"
fi

# Create .env.example file
echo "📝 Creating environment template..."
cat > .env.example << 'EOF'
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
EOF

echo ""
echo "🎉 Setup complete! Next steps:"
echo ""
echo "1. Review the learning path: docs/learning-path.md"
echo "2. Check Git workflow: docs/git-flow-strategy.md"
echo "3. Create develop branch:"
echo "   git checkout -b develop"
echo "   git push -u origin develop"
echo ""
echo "4. Start with Module 1:"
echo "   git checkout -b feature/module-1-foundations"
echo "   cd modules/module-01"
echo ""
echo "Happy learning! 🚀"
