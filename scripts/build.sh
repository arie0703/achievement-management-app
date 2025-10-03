#!/bin/bash

# Achievement Management Application Build Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="achievement-app"
API_NAME="achievement-api"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH=${COMMIT_HASH:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}

# Directories
BUILD_DIR="build"
DIST_DIR="dist"
SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPTS_DIR")"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    if ! command -v git &> /dev/null; then
        log_warning "Git is not installed - version info may be incomplete"
    fi
    
    log_success "Prerequisites check passed"
}

# Clean previous builds
clean_build() {
    log_info "Cleaning previous builds..."
    rm -rf "$ROOT_DIR/$BUILD_DIR"
    rm -rf "$ROOT_DIR/$DIST_DIR"
    log_success "Clean completed"
}

# Download dependencies
download_deps() {
    log_info "Downloading dependencies..."
    cd "$ROOT_DIR"
    go mod download
    go mod tidy
    log_success "Dependencies downloaded"
}

# Run tests
run_tests() {
    log_info "Running tests..."
    cd "$ROOT_DIR"
    
    if ! go test -v ./...; then
        log_error "Tests failed"
        exit 1
    fi
    
    log_success "All tests passed"
}

# Build for single platform
build_platform() {
    local os=$1
    local arch=$2
    local output_dir="$ROOT_DIR/$BUILD_DIR/$os-$arch"
    
    log_info "Building for $os/$arch..."
    
    mkdir -p "$output_dir"
    
    local ldflags="-ldflags=-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.CommitHash=$COMMIT_HASH"
    
    # Build API server
    if [ "$os" = "windows" ]; then
        GOOS=$os GOARCH=$arch go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.CommitHash=$COMMIT_HASH" -o "$output_dir/${API_NAME}.exe" ./cmd/api
        GOOS=$os GOARCH=$arch go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.CommitHash=$COMMIT_HASH" -o "$output_dir/${APP_NAME}.exe" ./cmd/cli
    else
        GOOS=$os GOARCH=$arch go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.CommitHash=$COMMIT_HASH" -o "$output_dir/$API_NAME" ./cmd/api
        GOOS=$os GOARCH=$arch go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.CommitHash=$COMMIT_HASH" -o "$output_dir/$APP_NAME" ./cmd/cli
    fi
    
    log_success "Build completed for $os/$arch"
}

# Build all platforms
build_all() {
    log_info "Building for all platforms..."
    
    cd "$ROOT_DIR"
    
    # Linux
    build_platform "linux" "amd64"
    build_platform "linux" "arm64"
    
    # macOS
    build_platform "darwin" "amd64"
    build_platform "darwin" "arm64"
    
    # Windows
    build_platform "windows" "amd64"
    
    log_success "All platform builds completed"
}

# Build without tests
build_only() {
    log_info "Building for all platforms (skipping tests)..."
    
    cd "$ROOT_DIR"
    
    # Linux
    build_platform "linux" "amd64"
    build_platform "linux" "arm64"
    
    # macOS
    build_platform "darwin" "amd64"
    build_platform "darwin" "arm64"
    
    # Windows
    build_platform "windows" "amd64"
    
    log_success "All platform builds completed"
}

# Create distribution packages
create_packages() {
    log_info "Creating distribution packages..."
    
    mkdir -p "$ROOT_DIR/$DIST_DIR"
    cd "$ROOT_DIR"
    
    # Linux amd64
    tar -czf "$DIST_DIR/${APP_NAME}-${VERSION}-linux-amd64.tar.gz" \
        -C "$BUILD_DIR/linux-amd64" "$API_NAME" "$APP_NAME" \
        -C "../../config" . \
        -C ".." "README.md" ".env.example"
    
    # Linux arm64
    tar -czf "$DIST_DIR/${APP_NAME}-${VERSION}-linux-arm64.tar.gz" \
        -C "$BUILD_DIR/linux-arm64" "$API_NAME" "$APP_NAME" \
        -C "../../config" . \
        -C ".." "README.md" ".env.example"
    
    # macOS amd64
    tar -czf "$DIST_DIR/${APP_NAME}-${VERSION}-darwin-amd64.tar.gz" \
        -C "$BUILD_DIR/darwin-amd64" "$API_NAME" "$APP_NAME" \
        -C "../../config" . \
        -C ".." "README.md" ".env.example"
    
    # macOS arm64
    tar -czf "$DIST_DIR/${APP_NAME}-${VERSION}-darwin-arm64.tar.gz" \
        -C "$BUILD_DIR/darwin-arm64" "$API_NAME" "$APP_NAME" \
        -C "../../config" . \
        -C ".." "README.md" ".env.example"
    
    # Windows
    cd "$BUILD_DIR/windows-amd64"
    zip -r "../../$DIST_DIR/${APP_NAME}-${VERSION}-windows-amd64.zip" \
        "${API_NAME}.exe" "${APP_NAME}.exe" "../../config/"* "../../README.md" "../../.env.example"
    cd "$ROOT_DIR"
    
    log_success "Distribution packages created"
}

# Generate checksums
generate_checksums() {
    log_info "Generating checksums..."
    
    cd "$ROOT_DIR/$DIST_DIR"
    
    if command -v sha256sum &> /dev/null; then
        sha256sum *.tar.gz *.zip > checksums.txt
    elif command -v shasum &> /dev/null; then
        shasum -a 256 *.tar.gz *.zip > checksums.txt
    else
        log_warning "No checksum utility found (sha256sum or shasum)"
        return
    fi
    
    log_success "Checksums generated"
}

# Print build info
print_build_info() {
    log_info "Build Information:"
    echo "  Version: $VERSION"
    echo "  Build Time: $BUILD_TIME"
    echo "  Commit Hash: $COMMIT_HASH"
    echo "  Go Version: $(go version)"
}

# Main execution
main() {
    local command=${1:-"all"}
    
    case $command in
        "clean")
            clean_build
            ;;
        "deps")
            download_deps
            ;;
        "test")
            run_tests
            ;;
        "build")
            check_prerequisites
            download_deps
            run_tests
            build_all
            ;;
        "build-only")
            check_prerequisites
            download_deps
            build_only
            ;;
        "package")
            check_prerequisites
            download_deps
            run_tests
            build_all
            create_packages
            generate_checksums
            ;;
        "all")
            print_build_info
            check_prerequisites
            clean_build
            download_deps
            run_tests
            build_all
            create_packages
            generate_checksums
            log_success "Build process completed successfully!"
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  clean      - Clean build artifacts"
            echo "  deps       - Download dependencies"
            echo "  test       - Run tests"
            echo "  build      - Build binaries for all platforms (with tests)"
            echo "  build-only - Build binaries for all platforms (skip tests)"
            echo "  package    - Build and create distribution packages"
            echo "  all        - Run complete build process (default)"
            echo "  help       - Show this help message"
            ;;
        *)
            log_error "Unknown command: $command"
            echo "Use '$0 help' for usage information"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"