# mcpipboy justfile - Build automation and development tasks

_default:
    just --list

# Build the application
build:
    @echo "Building mcpipboy..."
    go build -o bin/mcpipboy ./cmd/mcpipboy

# Build release binaries with static linking
build-release:
    @echo "Building release binaries..."
    @mkdir -p dist
    # Linux
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/mcpipboy-linux-amd64 ./cmd/mcpipboy
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/mcpipboy-linux-arm64 ./cmd/mcpipboy
    # macOS
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/mcpipboy-darwin-amd64 ./cmd/mcpipboy
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/mcpipboy-darwin-arm64 ./cmd/mcpipboy
    # Windows
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/mcpipboy-windows-amd64.exe ./cmd/mcpipboy
    @echo "Release binaries built in dist/"

# Create GitHub release
release:
    @echo "Creating GitHub release..."
    @current_version=$$(cat VERSION) && \
    echo "Creating release v$$current_version" && \
    gh release create "v$$current_version" \
        --title "mcpipboy v$$current_version" \
        --notes "Release v$$current_version of mcpipboy - MCP server for AI agents" \
        dist/mcpipboy-linux-amd64 \
        dist/mcpipboy-linux-arm64 \
        dist/mcpipboy-darwin-amd64 \
        dist/mcpipboy-darwin-arm64 \
        dist/mcpipboy-windows-amd64.exe

# Generate release notes from git commits since last tag
release-notes:
    @echo "Generating release notes..."
    @git log --oneline --pretty=format:"- %s" --since="1 month ago"

# Build and release (builds binaries then creates GitHub release)
build-and-release: build-release release
    @echo "Build and release complete!"

# Run tests
test:
    @echo "Running tests..."
    go test -v ./...

# Run tests with coverage
test-coverage:
    @echo "Running tests with coverage..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

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

# Development mode - build and run
dev: build
    @echo "Running mcpipboy in development mode..."
    ./bin/mcpipboy

# Install dependencies
deps:
    @echo "Installing dependencies..."
    go mod download
    go mod tidy

# Install binary to $GOPATH/bin
install: build
    @echo "Installing mcpipboy..."
    cp bin/mcpipboy $(shell go env GOPATH)/bin/

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf bin/ dist/
    rm -f coverage.out coverage.html

# Pre-commit validation
check: fmt vet test
    @echo "All checks passed!"

# Get current version from VERSION file
get-version:
    @cat VERSION

# Set version in VERSION file
set-version version:
    @echo "{{version}}" > VERSION
    @echo "Version set to {{version}}"

# Bump patch version (0.1.0 -> 0.1.1)
bump-patch:
    @current=$(cat VERSION) && \
    major=$(echo $current | cut -d. -f1) && \
    minor=$(echo $current | cut -d. -f2) && \
    patch=$(echo $current | cut -d. -f3) && \
    new_version="$major.$minor.$((patch + 1))" && \
    echo "$new_version" > VERSION && \
    echo "Version bumped to $new_version"

# Bump minor version (0.1.0 -> 0.2.0)
bump-minor:
    @current=$(cat VERSION) && \
    major=$(echo $current | cut -d. -f1) && \
    minor=$(echo $current | cut -d. -f2) && \
    new_version="$major.$((minor + 1)).0" && \
    echo "$new_version" > VERSION && \
    echo "Version bumped to $new_version"

# Bump major version (0.1.0 -> 1.0.0)
bump-major:
    @current=$(cat VERSION) && \
    major=$(echo $current | cut -d. -f1) && \
    new_version="$((major + 1)).0.0" && \
    echo "$new_version" > VERSION && \
    echo "Version bumped to $new_version"

# Bump prerelease version (0.1.0 -> 0.1.0-alpha1)
bump-prerelease:
    @current=$(cat VERSION) && \
    if echo $current | grep -q "-"; then \
        base=$(echo $current | cut -d- -f1) && \
        prerelease=$(echo $current | cut -d- -f2) && \
        if echo $prerelease | grep -q "alpha"; then \
            alpha_num=$(echo $prerelease | sed 's/alpha//') && \
            if [ -z "$alpha_num" ]; then alpha_num=1; else alpha_num=$((alpha_num + 1)); fi && \
            new_version="$base-alpha$alpha_num"; \
        else \
            new_version="$base-alpha1"; \
        fi; \
    else \
        new_version="$current-alpha1"; \
    fi && \
    echo "$new_version" > VERSION && \
    echo "Version bumped to $new_version"

