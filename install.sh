#!/bin/bash

# GitCleaner Installation Script
# This script installs gitcleaner either from local source or GitHub releases
# Usage:
#   Local:  ./install.sh
#   Remote: curl -sSL https://raw.githubusercontent.com/theburrowhub/gitcleaner/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="gitcleaner"
REPO_OWNER="theburrowhub" # TODO: Update with your GitHub username
REPO_NAME="gitcleaner"
INSTALL_DIR="${HOME}/.local/bin"
BACKUP_DIR="${HOME}/.local/backup"

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
    exit 1
}

detect_installation_mode() {
    # Check if we're running from a local git repository with gitcleaner source
    if [ -f "go.mod" ] && [ -f "main.go" ] && [ -d ".git" ]; then
        # Additional check: verify it's actually the gitcleaner repository
        if grep -q "module gitcleaner" go.mod 2>/dev/null; then
            print_info "Detected local gitcleaner repository - will build from source"
            return 0  # Local mode
        fi
    fi
    
    print_info "Detected remote execution - will download from GitHub releases"
    return 1  # Remote mode
}

get_latest_release() {
    local api_url="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
    local release_info
    
    if command -v curl &> /dev/null; then
        release_info=$(curl -s "$api_url" 2>/dev/null)
    elif command -v wget &> /dev/null; then
        release_info=$(wget -qO- "$api_url" 2>/dev/null)
    else
        print_error "Neither curl nor wget is available. Please install one of them."
    fi
    
    local tag_name=$(echo "$release_info" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/' 2>/dev/null)
    
    if [ -z "$tag_name" ]; then
        print_error "Failed to get latest release information"
    fi
    
    echo "$tag_name"
}

detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        linux*)
            os="linux"
            ;;
        darwin*)
            os="macos"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            ;;
    esac
    
    case "$arch" in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            ;;
    esac
    
    echo "${os}-${arch}"
}

download_and_install_remote() {
    print_info "Fetching latest release information..."
    local version=$(get_latest_release)
    local platform=$(detect_platform)
    local archive_name="${BINARY_NAME}-${platform}.tar.gz"
    local download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${version}/${archive_name}"
    local temp_dir=$(mktemp -d)
    
    print_info "Downloading ${BINARY_NAME} ${version} for ${platform}..."
    print_info "Download URL: ${download_url}"
    
    # Download the archive
    if command -v curl &> /dev/null; then
        curl -sSL "$download_url" -o "${temp_dir}/${archive_name}" || print_error "Failed to download ${archive_name}"
    elif command -v wget &> /dev/null; then
        wget -q "$download_url" -O "${temp_dir}/${archive_name}" || print_error "Failed to download ${archive_name}"
    else
        print_error "Neither curl nor wget is available. Please install one of them."
    fi
    
    print_info "Extracting archive..."
    cd "$temp_dir"
    tar -xzf "$archive_name" || print_error "Failed to extract archive"
    
    # Find the binary (it might be named differently in the archive)
    local binary_path=$(find . -name "${BINARY_NAME}-${platform}" -type f | head -n1)
    if [ -z "$binary_path" ]; then
        print_error "Binary not found in archive"
    fi
    
    print_info "Installing binary to ${INSTALL_DIR}..."
    mkdir -p "$INSTALL_DIR"
    cp "$binary_path" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Cleanup
    rm -rf "$temp_dir"
    
    print_success "Downloaded and installed ${BINARY_NAME} ${version}"
}

check_go_dependencies() {
    print_info "Checking Go dependencies..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.19 or later."
    fi
    
    # Check Go version
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | grep -oE '[0-9]+\.[0-9]+')
    REQUIRED_VERSION="1.19"
    
    if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
        print_error "Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or later."
    fi
    
    print_success "Go dependencies check passed (Go $GO_VERSION)"
}

build_and_install_local() {
    check_go_dependencies
    
    # Verify we're in the correct directory
    if [ ! -f "go.mod" ] || ! grep -q "module gitcleaner" go.mod 2>/dev/null; then
        print_error "Not in gitcleaner repository directory. Please run from the gitcleaner source directory."
    fi
    
    print_info "Building from local source..."
    
    # Get version information
    local version="dev"
    local commit_hash="unknown"
    local build_time=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
    
    # Try to get git information if available
    if command -v git &> /dev/null && git rev-parse --git-dir > /dev/null 2>&1; then
        version=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
        commit_hash=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    fi
    
    # Build flags - gitcleaner doesn't have cmd package, so we'll use main package
    local ldflags="-X main.Version=${version} -X main.CommitHash=${commit_hash} -X main.BuildTime=${build_time}"
    
    print_info "Version: $version"
    print_info "Commit: $commit_hash"
    print_info "Build time: $build_time"
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Build the binary
    print_info "Compiling binary..."
    if ! go build -ldflags "$ldflags" -o "${INSTALL_DIR}/${BINARY_NAME}" .; then
        print_error "Failed to build binary. Make sure you're in the gitcleaner repository directory and Go modules are working."
    fi
    
    # Make it executable
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    print_success "Built and installed from local source"
}

backup_existing() {
    local existing_binary="$INSTALL_DIR/$BINARY_NAME"
    
    if [ -f "$existing_binary" ]; then
        print_info "Backing up existing installation..."
        
        # Get current version if possible
        local current_version="unknown"
        if [ -x "$existing_binary" ]; then
            current_version=$("$existing_binary" --version 2>/dev/null | head -n1 | grep -oE 'GitCleaner [^ ]+' | cut -d' ' -f2 || echo "unknown")
        fi
        
        mkdir -p "$BACKUP_DIR"
        local backup_file="$BACKUP_DIR/${BINARY_NAME}-${current_version}-$(date +%Y%m%d-%H%M%S)"
        cp "$existing_binary" "$backup_file"
        
        print_success "Backed up to: $backup_file"
        return 0
    fi
    
    print_info "No existing installation found"
    return 1
}

verify_installation() {
    local installed_binary="$INSTALL_DIR/$BINARY_NAME"
    
    print_info "Verifying installation..."
    
    if [ ! -f "$installed_binary" ]; then
        print_error "Installation failed: binary not found"
    fi
    
    if [ ! -x "$installed_binary" ]; then
        print_error "Installation failed: binary is not executable"
    fi
    
    # Test the binary
    local version_output
    if version_output=$("$installed_binary" --version 2>&1); then
        print_success "Installation verified"
        echo "$version_output"
    else
        print_error "Installation failed: binary is not working correctly"
    fi
}

update_path() {
    print_info "Checking PATH configuration..."
    
    # Check if install directory is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        print_warning "Installation directory is not in PATH"
        print_info "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo ""
        echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
        print_info "Then restart your shell or run: source ~/.bashrc (or ~/.zshrc)"
        echo ""
        print_info "Alternatively, you can run gitcleaner directly: $INSTALL_DIR/$BINARY_NAME"
    else
        print_success "Installation directory is already in PATH"
    fi
}

show_usage() {
    cat << EOF
GitCleaner Installation Script

This script automatically detects whether it's running locally or remotely:
- Local:  Builds from source code (requires Go)
- Remote: Downloads latest release from GitHub

Usage: 
    ./install.sh [OPTIONS]                    # Local installation
    curl -sSL <raw-url>/install.sh | bash     # Remote installation

Options:
    -h, --help          Show this help message
    -d, --dir DIR       Set installation directory (default: ~/.local/bin)
    --force             Force reinstallation without backup
    --local             Force local build mode
    --remote            Force remote download mode

Examples:
    ./install.sh                              # Auto-detect mode
    ./install.sh -d /usr/local/bin           # Install to custom directory
    ./install.sh --force --local            # Force local build
    curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/gitcleaner/main/install.sh | bash

EOF
}

# Parse command line arguments
FORCE_INSTALL=false
FORCE_LOCAL=false
FORCE_REMOTE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -d|--dir)
            INSTALL_DIR="$2"
            BACKUP_DIR="${INSTALL_DIR}/../backup"
            shift 2
            ;;
        --force)
            FORCE_INSTALL=true
            shift
            ;;
        --local)
            FORCE_LOCAL=true
            shift
            ;;
        --remote)
            FORCE_REMOTE=true
            shift
            ;;
        *)
            print_error "Unknown option: $1. Use --help for usage information."
            ;;
    esac
done

# Main installation process
main() {
    echo "🧹 GitCleaner Installation Script"
    echo "================================"
    echo ""
    
    print_info "Installation directory: $INSTALL_DIR"
    print_info "Backup directory: $BACKUP_DIR"
    echo ""
    
    # Determine installation mode
    local is_local_mode=false
    
    if [ "$FORCE_LOCAL" = true ]; then
        is_local_mode=true
        print_info "Forced local build mode"
    elif [ "$FORCE_REMOTE" = true ]; then
        is_local_mode=false
        print_info "Forced remote download mode"
    else
        if detect_installation_mode; then
            is_local_mode=true
        else
            is_local_mode=false
        fi
    fi
    
    # Handle existing installation
    if [ "$FORCE_INSTALL" = false ]; then
        backup_existing || true  # Don't fail if no existing installation
    fi
    
    # Install based on mode
    if [ "$is_local_mode" = true ]; then
        build_and_install_local
    else
        download_and_install_remote
    fi
    
    verify_installation
    update_path
    
    echo ""
    print_success "Installation completed successfully!"
    echo ""
    print_info "Next steps:"
    echo "  1. Navigate to any git repository"
    echo "  2. Run: gitcleaner"
    echo "  3. Select branches to delete using Tab/Space"
    echo "  4. Press Enter to confirm deletion"
    echo ""
    print_info "For more information, run: gitcleaner --help"
}

# Run main function
main "$@"