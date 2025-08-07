#!/bin/bash

# cons-coder installation script
# Used to compile and install cons-coder to system PATH

set -e  # Exit immediately on error

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or higher."
        echo "Visit https://golang.org/dl/ to download and install"
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    MIN_VERSION="1.21"
    
    if [ "$(printf '%s\n' "$MIN_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
        print_error "Go version is too old. Required: Go $MIN_VERSION or higher, Current: $GO_VERSION"
        exit 1
    fi
    
    print_info "Detected Go version: $GO_VERSION"
}

# Check GOPATH
check_gopath() {
    if [ -z "$GOPATH" ]; then
        GOPATH=$(go env GOPATH)
        if [ -z "$GOPATH" ]; then
            print_error "Unable to determine GOPATH"
            exit 1
        fi
    fi
    
    GOBIN="$GOPATH/bin"
    if [ ! -d "$GOBIN" ]; then
        print_info "Creating $GOBIN directory..."
        mkdir -p "$GOBIN"
    fi
    
    print_info "GOPATH: $GOPATH"
    print_info "GOBIN: $GOBIN"
}

# Download dependencies
download_deps() {
    print_info "Downloading project dependencies..."
    go mod download
    print_success "Dependencies downloaded successfully"
}

# Run tests (optional)
run_tests() {
    if [ "$1" == "--with-tests" ] || [ "$1" == "-t" ]; then
        print_info "Running tests..."
        go test ./...
        print_success "All tests passed"
    fi
}

# Compile and install
install_binary() {
    print_info "Compiling and installing cons-coder..."
    
    # Get version information
    VERSION="1.0.0"
    BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    # Set build flags
    LDFLAGS="-X 'main.Version=$VERSION' -X 'main.BuildTime=$BUILD_TIME' -X 'main.GitCommit=$GIT_COMMIT'"
    
    # Compile and install
    go install -ldflags "$LDFLAGS" .
    
    print_success "cons-coder has been installed to $GOBIN"
}

# Check PATH
check_path() {
    if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
        print_warning "$GOBIN is not in PATH"
        echo ""
        echo "Please add the following line to your shell configuration file (~/.bashrc, ~/.zshrc, etc):"
        echo ""
        echo "  export PATH=\"\$PATH:$GOBIN\""
        echo ""
        echo "Then run:"
        echo "  source ~/.bashrc  # or source ~/.zshrc"
        echo ""
    else
        print_success "PATH is configured correctly"
    fi
}

# Verify installation
verify_installation() {
    if command -v cons-coder &> /dev/null; then
        print_success "cons-coder installed successfully!"
        echo ""
        cons-coder --version 2>/dev/null || echo "cons-coder version $VERSION"
        echo ""
        echo "Usage:"
        echo "  cons-coder -d <XML_DIR> -o <OUTPUT_DIR> -l <LANGUAGE>"
        echo ""
        echo "Examples:"
        echo "  cons-coder -d ./data -o ./output/python -l python"
        echo "  cons-coder -d ./data -o ./output/swift -l swift"
        echo ""
        echo "Supported languages: python, go, java, swift, kotlin, typescript, javascript"
    else
        print_error "Installation verification failed, please check PATH settings"
        exit 1
    fi
}

# Uninstall function
uninstall() {
    print_info "Uninstalling cons-coder..."
    
    BINARY_PATH="$GOBIN/cons-coder"
    if [ -f "$BINARY_PATH" ]; then
        rm -f "$BINARY_PATH"
        print_success "cons-coder has been uninstalled from $GOBIN"
    else
        print_warning "cons-coder installation not found"
    fi
}

# Main function
main() {
    echo "======================================"
    echo "   cons-coder Installation Script"
    echo "======================================"
    echo ""
    
    # Handle command line arguments
    case "$1" in
        uninstall|--uninstall|-u)
            uninstall
            exit 0
            ;;
        help|--help|-h)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --with-tests, -t    Run tests before installation"
            echo "  --uninstall, -u     Uninstall cons-coder"
            echo "  --help, -h          Show this help message"
            echo ""
            exit 0
            ;;
    esac
    
    # Installation workflow
    check_go
    check_gopath
    download_deps
    run_tests "$1"
    install_binary
    check_path
    verify_installation
}

# Run main function
main "$@"