#!/bin/bash

# LogAid Installation Script
# This script downloads and installs the latest version of LogAid

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="ayushsharma-1/LogAid"
BINARY_NAME="logaid"
INSTALL_DIR="/usr/local/bin"

# Utility functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and Architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $os in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    log_info "Detected platform: $OS-$ARCH"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command_exists curl && ! command_exists wget; then
        log_error "Either curl or wget is required to download LogAid"
        exit 1
    fi
    
    if ! command_exists tar; then
        log_error "tar is required to extract LogAid"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Get latest release version
get_latest_version() {
    log_info "Fetching latest release information..."
    
    if command_exists curl; then
        LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command_exists wget; then
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    if [ -z "$LATEST_VERSION" ]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi
    
    log_info "Latest version: $LATEST_VERSION"
}

# Download and install LogAid
install_logaid() {
    local download_url="https://github.com/$REPO/releases/download/$LATEST_VERSION/logaid-$OS-$ARCH.tar.gz"
    local temp_dir=$(mktemp -d)
    local temp_file="$temp_dir/logaid.tar.gz"
    
    log_info "Downloading LogAid $LATEST_VERSION..."
    log_info "Download URL: $download_url"
    
    # Download the archive
    if command_exists curl; then
        curl -L -o "$temp_file" "$download_url"
    elif command_exists wget; then
        wget -O "$temp_file" "$download_url"
    fi
    
    if [ ! -f "$temp_file" ]; then
        log_error "Failed to download LogAid"
        exit 1
    fi
    
    log_info "Extracting archive..."
    tar -xzf "$temp_file" -C "$temp_dir"
    
    # Find the binary (it might be named differently)
    local binary_path=""
    if [ -f "$temp_dir/logaid-$OS-$ARCH" ]; then
        binary_path="$temp_dir/logaid-$OS-$ARCH"
    elif [ -f "$temp_dir/logaid" ]; then
        binary_path="$temp_dir/logaid"
    else
        log_error "Could not find LogAid binary in archive"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$binary_path"
    
    # Install binary
    log_info "Installing LogAid to $INSTALL_DIR..."
    
    if [ -w "$INSTALL_DIR" ]; then
        cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    else
        log_info "Installing with sudo (requires administrator privileges)..."
        sudo cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
    
    log_success "LogAid installed successfully!"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    if command_exists "$BINARY_NAME"; then
        local installed_version=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        log_success "LogAid is installed and working!"
        log_info "Installed version: $installed_version"
        log_info "Installation path: $(which $BINARY_NAME)"
    else
        log_error "LogAid installation verification failed"
        log_info "You may need to add $INSTALL_DIR to your PATH"
        exit 1
    fi
}

# Show usage information
show_usage() {
    log_info "LogAid has been installed successfully!"
    echo
    echo "Usage examples:"
    echo "  $BINARY_NAME --help                 # Show help"
    echo "  $BINARY_NAME --version              # Show version"
    echo "  git clone invalid-repo 2>&1 | $BINARY_NAME  # Pipe error output"
    echo "  $BINARY_NAME 'apt install nonexistent'      # Direct command analysis"
    echo
    echo "For more information, visit: https://github.com/$REPO"
}

# Handle installation options
handle_options() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                if [ -n "$2" ]; then
                    LATEST_VERSION="$2"
                    log_info "Installing specific version: $LATEST_VERSION"
                    shift
                fi
                shift
                ;;
            --install-dir)
                if [ -n "$2" ]; then
                    INSTALL_DIR="$2"
                    log_info "Installing to custom directory: $INSTALL_DIR"
                    shift
                fi
                shift
                ;;
            --help|-h)
                echo "LogAid Installation Script"
                echo
                echo "Usage: $0 [OPTIONS]"
                echo
                echo "Options:"
                echo "  --version VERSION     Install specific version"
                echo "  --install-dir DIR     Install to custom directory (default: /usr/local/bin)"
                echo "  --help, -h            Show this help message"
                echo
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
}

# Main installation process
main() {
    echo "╔══════════════════════════════════════╗"
    echo "║            LogAid Installer          ║"
    echo "║    Intelligent CLI Error Assistant   ║"
    echo "╚══════════════════════════════════════╝"
    echo
    
    # Handle command line options
    handle_options "$@"
    
    # Detect platform
    detect_platform
    
    # Check prerequisites
    check_prerequisites
    
    # Get latest version if not specified
    if [ -z "$LATEST_VERSION" ]; then
        get_latest_version
    fi
    
    # Install LogAid
    install_logaid
    
    # Verify installation
    verify_installation
    
    # Show usage information
    show_usage
    
    log_success "Installation completed successfully! 🎉"
}

# Run main function with all arguments
main "$@"
