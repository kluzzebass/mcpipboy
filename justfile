# mcpipboy justfile - Build automation and development tasks

_default:
    just --list

# Build the application
build:
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
    DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    LDFLAGS="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"
    echo "Building mcpipboy $VERSION..."
    mkdir -p dist
    CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o dist/mcpipboy ./cmd/mcpipboy
    echo "Built dist/mcpipboy"

# Build release binaries with static linking for all platforms
build-all:
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
    DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    LDFLAGS="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"
    echo "Building mcpipboy $VERSION for all platforms..."
    mkdir -p dist
    # Linux
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-linux-amd64 ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-linux-arm64 ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-linux-386 ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-linux-arm ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-linux-riscv64 ./cmd/mcpipboy
    # macOS
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-darwin-amd64 ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-darwin-arm64 ./cmd/mcpipboy
    # Windows
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-windows-amd64.exe ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-windows-arm64.exe ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-windows-386.exe ./cmd/mcpipboy
    # FreeBSD
    CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-freebsd-amd64 ./cmd/mcpipboy
    CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/mcpipboy-freebsd-arm64 ./cmd/mcpipboy
    echo "Built binaries for all platforms in dist/"

# Run tests
test:
    @echo "Running tests..."
    @go test -v ./...

# Run tests with coverage and show percentages
test-coverage:
    @echo "Running tests with coverage..."
    go test -short -coverprofile=coverage.out ./... || true
    @if [ -f coverage.out ]; then \
        go tool cover -func=coverage.out | tail -1; \
    fi

# Generate HTML coverage report
coverage-html:
    @echo "Generating HTML coverage report..."
    @if [ -f coverage.out ]; then \
        go tool cover -html=coverage.out -o coverage.html && \
        echo "Coverage report generated: coverage.html"; \
    else \
        echo "Error: coverage.out not found. Run 'just test-coverage' first."; \
        exit 1; \
    fi

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...

# Check for common Go mistakes and suspicious constructs
vet:
    @echo "Running go vet..."
    go vet ./...

# Run golangci-lint (if installed)
lint:
    @echo "Running linter..."
    @if command -v golangci-lint >/dev/null 2>&1; then \
        golangci-lint run; \
    else \
        echo "golangci-lint not installed, skipping..."; \
    fi

# Development mode - build and test
dev:
    @just build
    @just test
    @echo "Development build complete!"

# Install dependencies
deps:
    @echo "Installing dependencies..."
    go mod download
    go mod tidy

# Install binary to $GOPATH/bin
install:
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
    DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    LDFLAGS="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"
    echo "Installing mcpipboy $VERSION..."
    go install -ldflags "$LDFLAGS" ./cmd/mcpipboy

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf bin/ dist/
    rm -f coverage.out coverage.html

# Pre-commit validation
check: fmt vet test
    @echo "All checks passed!"

# Show project info
info:
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    echo "mcpipboy - MCP Server for AI Agent Tools"
    echo "========================================="
    echo "Version: $VERSION"
    echo "Go version: $(go version)"
    echo "Git commit: $(git rev-parse --short HEAD 2>/dev/null || echo 'none')"
    echo "Git branch: $(git branch --show-current 2>/dev/null || echo 'none')"
    echo "Build date: $(date)"

# Show current git tags
git-tags:
    @echo "Git tags:"
    @git tag --sort=-version:refname
