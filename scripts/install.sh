#!/bin/bash

# Achievement Management Application Installation Script

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
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/achievement-management"
SERVICE_DIR="/etc/systemd/system"

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

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_error "This script should not be run as root"
        log_info "Please run as a regular user with sudo privileges"
        exit 1
    fi
}

# Detect platform
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $arch in
        x86_64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    case $os in
        linux)
            PLATFORM="linux-$arch"
            ;;
        darwin)
            PLATFORM="darwin-$arch"
            INSTALL_DIR="/usr/local/bin"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    log_info "Detected platform: $PLATFORM"
}

# Download and extract binaries
install_binaries() {
    local version=${1:-"latest"}
    local temp_dir=$(mktemp -d)
    
    log_info "Installing binaries for version: $version"
    log_info "Installation directory: $INSTALL_DIR"
    
    # Create installation directory if it doesn't exist
    sudo mkdir -p "$INSTALL_DIR"
    
    # For this example, we'll copy from the build directory
    # In a real deployment, you would download from a release URL
    if [[ -f "build/$PLATFORM/$API_NAME" ]]; then
        log_info "Installing from local build..."
        sudo cp "build/$PLATFORM/$API_NAME" "$INSTALL_DIR/"
        sudo cp "build/$PLATFORM/$APP_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$API_NAME"
        sudo chmod +x "$INSTALL_DIR/$APP_NAME"
    else
        log_error "Binaries not found in build/$PLATFORM/"
        log_info "Please run 'make build-all' first"
        exit 1
    fi
    
    log_success "Binaries installed successfully"
}

# Install configuration files
install_config() {
    log_info "Installing configuration files..."
    
    sudo mkdir -p "$CONFIG_DIR"
    
    if [[ -d "config" ]]; then
        sudo cp -r config/* "$CONFIG_DIR/"
        sudo chmod 644 "$CONFIG_DIR"/*
    fi
    
    if [[ -f ".env.example" ]]; then
        sudo cp .env.example "$CONFIG_DIR/environment.example"
    fi
    
    log_success "Configuration files installed"
}

# Create systemd service (Linux only)
create_service() {
    if [[ "$PLATFORM" != linux-* ]]; then
        log_info "Skipping systemd service creation (not Linux)"
        return
    fi
    
    log_info "Creating systemd service..."
    
    cat << EOF | sudo tee "$SERVICE_DIR/achievement-api.service" > /dev/null
[Unit]
Description=Achievement Management API Server
After=network.target

[Service]
Type=simple
User=achievement
Group=achievement
WorkingDirectory=$CONFIG_DIR
ExecStart=$INSTALL_DIR/$API_NAME
Restart=always
RestartSec=5
Environment=ENVIRONMENT=production
Environment=CONFIG_PATH=$CONFIG_DIR

[Install]
WantedBy=multi-user.target
EOF
    
    # Create user for the service
    if ! id "achievement" &>/dev/null; then
        sudo useradd -r -s /bin/false achievement
        log_info "Created achievement user"
    fi
    
    # Set ownership
    sudo chown -R achievement:achievement "$CONFIG_DIR"
    
    # Reload systemd and enable service
    sudo systemctl daemon-reload
    sudo systemctl enable achievement-api.service
    
    log_success "Systemd service created and enabled"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    if command -v "$APP_NAME" &> /dev/null; then
        local version=$($APP_NAME --version)
        log_success "CLI tool installed: $version"
    else
        log_error "CLI tool not found in PATH"
    fi
    
    if command -v "$API_NAME" &> /dev/null; then
        log_success "API server installed"
    else
        log_error "API server not found in PATH"
    fi
}

# Show post-installation instructions
show_instructions() {
    log_success "Installation completed!"
    echo
    echo "Next steps:"
    echo "1. Configure your environment variables in $CONFIG_DIR/"
    echo "2. Set up your AWS credentials and DynamoDB tables"
    echo
    echo "Usage:"
    echo "  CLI: $APP_NAME --help"
    echo "  API: $API_NAME"
    echo
    
    if [[ "$PLATFORM" == linux-* ]]; then
        echo "Service management:"
        echo "  Start:   sudo systemctl start achievement-api"
        echo "  Stop:    sudo systemctl stop achievement-api"
        echo "  Status:  sudo systemctl status achievement-api"
        echo "  Logs:    sudo journalctl -u achievement-api -f"
        echo
    fi
}

# Uninstall function
uninstall() {
    log_info "Uninstalling Achievement Management Application..."
    
    # Stop and disable service (Linux only)
    if [[ "$PLATFORM" == linux-* ]] && [[ -f "$SERVICE_DIR/achievement-api.service" ]]; then
        sudo systemctl stop achievement-api.service || true
        sudo systemctl disable achievement-api.service || true
        sudo rm -f "$SERVICE_DIR/achievement-api.service"
        sudo systemctl daemon-reload
        log_info "Systemd service removed"
    fi
    
    # Remove binaries
    sudo rm -f "$INSTALL_DIR/$API_NAME"
    sudo rm -f "$INSTALL_DIR/$APP_NAME"
    log_info "Binaries removed"
    
    # Remove configuration (with confirmation)
    if [[ -d "$CONFIG_DIR" ]]; then
        read -p "Remove configuration directory $CONFIG_DIR? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            sudo rm -rf "$CONFIG_DIR"
            log_info "Configuration removed"
        fi
    fi
    
    # Remove user (Linux only)
    if [[ "$PLATFORM" == linux-* ]] && id "achievement" &>/dev/null; then
        read -p "Remove achievement user? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            sudo userdel achievement || true
            log_info "User removed"
        fi
    fi
    
    log_success "Uninstallation completed"
}

# Main execution
main() {
    local command=${1:-"install"}
    local version=${2:-"latest"}
    
    case $command in
        "install")
            check_root
            detect_platform
            install_binaries "$version"
            install_config
            create_service
            verify_installation
            show_instructions
            ;;
        "uninstall")
            detect_platform
            uninstall
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [command] [version]"
            echo ""
            echo "Commands:"
            echo "  install    - Install the application (default)"
            echo "  uninstall  - Uninstall the application"
            echo "  help       - Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 install latest"
            echo "  $0 uninstall"
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