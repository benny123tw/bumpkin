# Bumpkin - Semantic Version Tagger CLI
# https://github.com/casey/just

# Default recipe - show available commands
default:
    @just --list

# Build the binary
build:
    go build -o bin/bumpkin ./cmd/bumpkin

# Build with version info
build-release version:
    go build -ldflags "-X github.com/benny123tw/bumpkin/internal/cli.AppVersion={{version}} -X github.com/benny123tw/bumpkin/internal/cli.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/bumpkin ./cmd/bumpkin

# Run the CLI (interactive mode)
run *args:
    go run ./cmd/bumpkin {{args}}

# Run all tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Run tests with coverage
test-cov:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report: coverage.html"

# Run linter
lint:
    golangci-lint run ./...

# Format code
fmt:
    golangci-lint fmt ./...
    go fmt ./...

# Run tests and linter
check: test lint

# Clean build artifacts
clean:
    rm -rf bin/
    rm -f coverage.out coverage.html

# Install binary to $GOPATH/bin
install:
    go install ./cmd/bumpkin

# Tidy go modules
tidy:
    go mod tidy

# Update dependencies
update:
    go get -u ./...
    go mod tidy

# Show current version
version:
    @go run ./cmd/bumpkin --show-version

# Run interactive mode
interactive:
    go run ./cmd/bumpkin

# Dry run patch bump
dry-patch:
    go run ./cmd/bumpkin --patch --dry-run --yes

# Dry run minor bump
dry-minor:
    go run ./cmd/bumpkin --minor --dry-run --yes

# Dry run major bump
dry-major:
    go run ./cmd/bumpkin --major --dry-run --yes

# Dry run conventional commit analysis
dry-conventional:
    go run ./cmd/bumpkin --conventional --dry-run --yes
