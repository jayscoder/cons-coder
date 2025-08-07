#!/bin/bash

# cons-coder build script
# Cross-platform build and release script

set -e  # Exit on error

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Project information
PROJECT_NAME="cons-coder"
VERSION="1.0.0"
BUILD_DIR="build"
DIST_DIR="dist"

# Get build information
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_STATE=$(if [ -z "$(git status --porcelain 2>/dev/null)" ]; then echo "clean"; else echo "dirty"; fi)

# LDFLAGS for version injection
LDFLAGS="-s -w \
    -X 'main.Version=${VERSION}' \
    -X 'main.BuildTime=${BUILD_TIME}' \
    -X 'main.GitCommit=${GIT_COMMIT}' \
    -X 'main.GitTag=${GIT_TAG}' \
    -X 'main.GitBranch=${GIT_BRANCH}' \
    -X 'main.GitState=${GIT_STATE}'"

# Platform configurations
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "windows/amd64"
    "windows/386"
    "windows/arm64"
)

# Print functions
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

print_step() {
    echo -e "${CYAN}==>${NC} $1"
}

# Check prerequisites
check_requirements() {
    print_step "Checking requirements..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go version: $GO_VERSION"
    
    if ! command -v git &> /dev/null; then
        print_warning "Git is not installed, version info will be limited"
    fi
}

# Clean build artifacts
clean() {
    print_step "Cleaning build artifacts..."
    rm -rf ${BUILD_DIR} ${DIST_DIR}
    rm -f ${PROJECT_NAME}
    print_success "Clean complete"
}

# Run tests
run_tests() {
    print_step "Running tests..."
    go test -v -race -cover ./...
    print_success "All tests passed"
}

# Run linters
run_lint() {
    print_step "Running linters..."
    
    # Check if golangci-lint is installed
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
        print_success "Lint checks passed"
    else
        print_warning "golangci-lint not installed, skipping lint checks"
        echo "Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s"
    fi
}

# Format code
format_code() {
    print_step "Formatting code..."
    go fmt ./...
    go mod tidy
    print_success "Code formatted"
}

# Build for current platform
build_local() {
    print_step "Building for current platform..."
    
    mkdir -p ${BUILD_DIR}
    
    OUTPUT="${BUILD_DIR}/${PROJECT_NAME}"
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
        OUTPUT="${OUTPUT}.exe"
    fi
    
    go build -ldflags "${LDFLAGS}" -o ${OUTPUT} .
    
    print_success "Build complete: ${OUTPUT}"
    
    # Show binary info
    if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "cygwin" && "$OSTYPE" != "win32" ]]; then
        file ${OUTPUT}
        ls -lh ${OUTPUT}
    fi
}

# Build for specific platform
build_platform() {
    GOOS=$1
    GOARCH=$2
    
    OUTPUT_NAME="${PROJECT_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    if [ "${GOOS}" == "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    OUTPUT_PATH="${BUILD_DIR}/${GOOS}-${GOARCH}/${PROJECT_NAME}"
    if [ "${GOOS}" == "windows" ]; then
        OUTPUT_PATH="${OUTPUT_PATH}.exe"
    fi
    
    print_info "Building for ${GOOS}/${GOARCH}..."
    
    mkdir -p "${BUILD_DIR}/${GOOS}-${GOARCH}"
    
    GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 \
        go build -ldflags "${LDFLAGS}" -o ${OUTPUT_PATH} .
    
    # Create archive
    mkdir -p ${DIST_DIR}
    if [ "${GOOS}" == "windows" ]; then
        # Create zip for Windows
        cd ${BUILD_DIR}/${GOOS}-${GOARCH}
        zip -q ../../${DIST_DIR}/${OUTPUT_NAME%.exe}.zip *
        cd ../..
    else
        # Create tar.gz for Unix-like systems
        cd ${BUILD_DIR}/${GOOS}-${GOARCH}
        tar -czf ../../${DIST_DIR}/${OUTPUT_NAME}.tar.gz *
        cd ../..
    fi
    
    print_success "Built ${GOOS}/${GOARCH}"
}

# Build for all platforms
build_all() {
    print_step "Building for all platforms..."
    
    clean
    mkdir -p ${BUILD_DIR} ${DIST_DIR}
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r GOOS GOARCH <<< "${platform}"
        build_platform ${GOOS} ${GOARCH}
    done
    
    print_success "All platforms built successfully"
    
    # Show distribution files
    print_step "Distribution files:"
    ls -lh ${DIST_DIR}/
}

# Build Docker image
build_docker() {
    print_step "Building Docker image..."
    
    # Create Dockerfile if not exists
    if [ ! -f Dockerfile ]; then
        cat > Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o cons-coder .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/cons-coder .

ENTRYPOINT ["./cons-coder"]
EOF
        print_info "Created Dockerfile"
    fi
    
    docker build -t ${PROJECT_NAME}:${VERSION} .
    docker tag ${PROJECT_NAME}:${VERSION} ${PROJECT_NAME}:latest
    
    print_success "Docker image built: ${PROJECT_NAME}:${VERSION}"
}

# Create release
create_release() {
    print_step "Creating release ${VERSION}..."
    
    # Build all platforms
    build_all
    
    # Create checksums
    print_info "Generating checksums..."
    cd ${DIST_DIR}
    shasum -a 256 * > SHA256SUMS
    cd ..
    
    # Create release notes
    cat > ${DIST_DIR}/RELEASE_NOTES.md << EOF
# Release ${VERSION}

## Build Information
- Version: ${VERSION}
- Build Time: ${BUILD_TIME}
- Git Commit: ${GIT_COMMIT}
- Git Branch: ${GIT_BRANCH}
- Git State: ${GIT_STATE}

## Supported Platforms
- macOS (Intel/Apple Silicon)
- Linux (amd64/arm64/386)
- Windows (amd64/386/arm64)

## Installation

### macOS/Linux
\`\`\`bash
# Download the appropriate binary for your platform
tar -xzf cons-coder-${VERSION}-<OS>-<ARCH>.tar.gz
chmod +x cons-coder
sudo mv cons-coder /usr/local/bin/
\`\`\`

### Windows
1. Download the appropriate .zip file
2. Extract cons-coder.exe
3. Add to your PATH

## Usage
\`\`\`bash
cons-coder -d <XML_DIR> -o <OUTPUT_DIR> -l <LANGUAGE>
\`\`\`

## Supported Languages
- Python
- Go
- Java
- Swift
- Kotlin
- TypeScript
- JavaScript
EOF
    
    print_success "Release ${VERSION} created in ${DIST_DIR}/"
    print_info "Release files:"
    ls -lh ${DIST_DIR}/
}

# Show version info
show_version() {
    echo "${PROJECT_NAME} build script"
    echo "Version: ${VERSION}"
    echo "Build Time: ${BUILD_TIME}"
    echo "Git Commit: ${GIT_COMMIT}"
    echo "Git Branch: ${GIT_BRANCH}"
    echo "Git State: ${GIT_STATE}"
}

# Show help
show_help() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  build       Build for current platform (default)"
    echo "  all         Build for all platforms"
    echo "  release     Create a release with all platforms"
    echo "  docker      Build Docker image"
    echo "  test        Run tests"
    echo "  lint        Run linters"
    echo "  fmt         Format code"
    echo "  clean       Clean build artifacts"
    echo "  version     Show version information"
    echo "  help        Show this help message"
    echo ""
    echo "Platform-specific builds:"
    echo "  darwin-amd64    Build for macOS (Intel)"
    echo "  darwin-arm64    Build for macOS (Apple Silicon)"
    echo "  linux-amd64     Build for Linux (amd64)"
    echo "  linux-arm64     Build for Linux (arm64)"
    echo "  windows-amd64   Build for Windows (amd64)"
    echo ""
    echo "Examples:"
    echo "  $0              # Build for current platform"
    echo "  $0 all          # Build for all platforms"
    echo "  $0 release      # Create a release"
    echo "  $0 test         # Run tests"
    echo "  $0 linux-amd64  # Build for Linux amd64"
}

# Main function
main() {
    case "${1:-build}" in
        build)
            check_requirements
            build_local
            ;;
        all)
            check_requirements
            build_all
            ;;
        release)
            check_requirements
            run_tests
            create_release
            ;;
        docker)
            check_requirements
            build_docker
            ;;
        test)
            check_requirements
            run_tests
            ;;
        lint)
            check_requirements
            run_lint
            ;;
        fmt|format)
            check_requirements
            format_code
            ;;
        clean)
            clean
            ;;
        version)
            show_version
            ;;
        help|--help|-h)
            show_help
            ;;
        darwin-amd64)
            check_requirements
            build_platform darwin amd64
            ;;
        darwin-arm64)
            check_requirements
            build_platform darwin arm64
            ;;
        linux-amd64)
            check_requirements
            build_platform linux amd64
            ;;
        linux-arm64)
            check_requirements
            build_platform linux arm64
            ;;
        windows-amd64)
            check_requirements
            build_platform windows amd64
            ;;
        *)
            print_error "Unknown command: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"