#!/bin/bash
set -e

# LogAid Installation Script
# This script installs LogAid on Linux and macOS systems

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="ayushsharma-1/LogAid"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="logaid"

# Functions
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

detect_os() {
    case "$OSTYPE" in
        linux*)   echo "linux" ;;
        darwin*)  echo "macos" ;;
        *)        echo "unsupported" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) echo "unsupported" ;;
    esac
}

check_dependencies() {
    print_info "Checking dependencies..."
    
    if ! command -v curl >/dev/null 2>&1; then
        print_error "curl is required but not installed."
        exit 1
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        print_error "tar is required but not installed."
        exit 1
    fi
    
    print_success "All dependencies are available"
}

get_latest_release() {
    print_info "Getting latest release information..."
    curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

download_and_install() {
    local os="$1"
    local arch="$2"
    local version="$3"
    
    local binary_name="${BINARY_NAME}-${os}-${arch}"
    local download_url="https://github.com/$REPO/releases/download/$version/$binary_name"
    
    print_info "Downloading LogAid $version for $os-$arch..."
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Download binary
    if ! curl -L -o "$binary_name" "$download_url"; then
        print_error "Failed to download LogAid from $download_url"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Make executable
    chmod +x "$binary_name"
    
    # Install to system
    print_info "Installing LogAid to $INSTALL_DIR..."
    
    if [[ $EUID -eq 0 ]]; then
        # Running as root
        mv "$binary_name" "$INSTALL_DIR/$BINARY_NAME"
    else
        # Need sudo
        sudo mv "$binary_name" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup
    cd - >/dev/null
    rm -rf "$temp_dir"
    
    print_success "LogAid installed successfully!"
}

setup_config() {
    print_info "Setting up configuration..."
    
    local config_dir="$HOME/.logaid"
    local env_file="$config_dir/.env"
    
    # Create config directory
    mkdir -p "$config_dir"
    mkdir -p "$config_dir/logs"
    mkdir -p "$config_dir/plugins"
    
    # Create default .env file if it doesn't exist
    if [[ ! -f "$env_file" ]]; then
        cat > "$env_file" << EOF
# LogAid Configuration
# Copy this file to .env and configure as needed

# AI Configuration
LOGAID_AI_PROVIDER=gemini
# LOGAID_API_KEY=your-api-key-here
LOGAID_AI_MODEL=gemini-1.5-flash
LOGAID_MAX_TOKENS=1000
LOGAID_TEMPERATURE=0.3

# Logging Configuration
LOGAID_LOG_LEVEL=info
LOGAID_LOG_PATH=$config_dir/logs/history.json

# Plugin Configuration
LOGAID_PLUGIN_DIR=$config_dir/plugins
LOGAID_ENABLED_PLUGINS=apt,git,npm,docker,kubernetes,generic

# Terminal Configuration
LOGAID_SHELL=$SHELL
LOGAID_PROMPT_TIMEOUT=30

# Feature Flags
LOGAID_ENABLE_LOCAL_FALLBACK=false
LOGAID_ENABLE_LOGGING=true
LOGAID_ENABLE_COLORS=true

# API Keys (set one of these)
# GEMINI_API_KEY=your-gemini-api-key
# OPENAI_API_KEY=your-openai-api-key
EOF
        print_success "Created default configuration at $env_file"
        print_warning "Don't forget to set your API key in $env_file"
    else
        print_info "Configuration file already exists at $env_file"
    fi
}

show_next_steps() {
    echo
    print_success "LogAid installation completed!"
    echo
    print_info "Next steps:"
    echo "  1. Set your API key:"
    echo "     export GEMINI_API_KEY=your-api-key"
    echo "     # or"
    echo "     export OPENAI_API_KEY=your-api-key"
    echo
    echo "  2. Test the installation:"
    echo "     logaid test"
    echo
    echo "  3. Start using LogAid:"
    echo "     logaid run"
    echo
    echo "  4. View help:"
    echo "     logaid --help"
    echo
    print_info "Configuration file: $HOME/.logaid/.env"
    print_info "Documentation: https://github.com/$REPO"
}

main() {
    echo "🚀 LogAid Installation Script"
    echo "============================="
    echo
    
    # Detect system
    local os=$(detect_os)
    local arch=$(detect_arch)
    
    if [[ "$os" == "unsupported" ]]; then
        print_error "Unsupported operating system: $OSTYPE"
        exit 1
    fi
    
    if [[ "$arch" == "unsupported" ]]; then
        print_error "Unsupported architecture: $(uname -m)"
        exit 1
    fi
    
    print_info "Detected system: $os-$arch"
    
    # Check dependencies
    check_dependencies
    
    # Get latest version
    local version=$(get_latest_release)
    if [[ -z "$version" ]]; then
        print_error "Failed to get latest release information"
        exit 1
    fi
    
    print_info "Latest version: $version"
    
    # Download and install
    download_and_install "$os" "$arch" "$version"
    
    # Setup configuration
    setup_config
    
    # Show next steps
    show_next_steps
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "LogAid Installation Script"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --uninstall    Uninstall LogAid"
        echo
        echo "This script will:"
        echo "  1. Detect your operating system and architecture"
        echo "  2. Download the latest LogAid release"
        echo "  3. Install it to $INSTALL_DIR"
        echo "  4. Set up default configuration"
        echo
        exit 0
        ;;
    --uninstall)
        print_info "Uninstalling LogAid..."
        if [[ -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
            if [[ $EUID -eq 0 ]]; then
                rm "$INSTALL_DIR/$BINARY_NAME"
            else
                sudo rm "$INSTALL_DIR/$BINARY_NAME"
            fi
            print_success "LogAid uninstalled successfully!"
        else
            print_warning "LogAid is not installed"
        fi
        exit 0
        ;;
    "")
        main
        ;;
    *)
        print_error "Unknown option: $1"
        echo "Use --help for usage information"
        exit 1
        ;;
esac
