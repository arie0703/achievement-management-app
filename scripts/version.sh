#!/bin/bash

# Version management script for Achievement Management Application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Get current version
get_current_version() {
    if git describe --tags --exact-match HEAD 2>/dev/null; then
        # On a tag
        git describe --tags --exact-match HEAD
    elif git describe --tags 2>/dev/null; then
        # After a tag
        git describe --tags
    else
        # No tags yet
        echo "v0.0.0-$(git rev-parse --short HEAD)"
    fi
}

# Get next version
get_next_version() {
    local current_version=$(get_current_version | sed 's/^v//' | sed 's/-.*//')
    local version_type=${1:-"patch"}
    
    IFS='.' read -ra VERSION_PARTS <<< "$current_version"
    local major=${VERSION_PARTS[0]:-0}
    local minor=${VERSION_PARTS[1]:-0}
    local patch=${VERSION_PARTS[2]:-0}
    
    case $version_type in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch")
            patch=$((patch + 1))
            ;;
        *)
            log_error "Invalid version type: $version_type"
            log_info "Valid types: major, minor, patch"
            exit 1
            ;;
    esac
    
    echo "v$major.$minor.$patch"
}

# Create a new tag
create_tag() {
    local version=$1
    local message=${2:-"Release $version"}
    
    if git tag -l | grep -q "^$version$"; then
        log_error "Tag $version already exists"
        exit 1
    fi
    
    log_info "Creating tag: $version"
    log_info "Message: $message"
    
    git tag -a "$version" -m "$message"
    log_success "Tag created: $version"
    
    log_info "To push the tag, run: git push origin $version"
}

# Show version information
show_version_info() {
    local current=$(get_current_version)
    local commit=$(git rev-parse --short HEAD)
    local branch=$(git branch --show-current)
    local build_time=$(date -u '+%Y-%m-%d_%H:%M:%S')
    
    echo "Version Information:"
    echo "  Current Version: $current"
    echo "  Commit Hash: $commit"
    echo "  Branch: $branch"
    echo "  Build Time: $build_time"
    echo
    echo "Next Version Options:"
    echo "  Patch: $(get_next_version patch)"
    echo "  Minor: $(get_next_version minor)"
    echo "  Major: $(get_next_version major)"
}

# List all tags
list_tags() {
    log_info "All tags:"
    git tag -l --sort=-version:refname | head -20
}

# Delete a tag
delete_tag() {
    local version=$1
    
    if [[ -z "$version" ]]; then
        log_error "Version not specified"
        exit 1
    fi
    
    if ! git tag -l | grep -q "^$version$"; then
        log_error "Tag $version does not exist"
        exit 1
    fi
    
    read -p "Delete tag $version? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git tag -d "$version"
        log_success "Tag deleted: $version"
        log_info "To delete from remote, run: git push origin :refs/tags/$version"
    else
        log_info "Tag deletion cancelled"
    fi
}

# Main execution
main() {
    local command=${1:-"info"}
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "Not in a git repository"
        exit 1
    fi
    
    case $command in
        "info"|"current")
            show_version_info
            ;;
        "list")
            list_tags
            ;;
        "tag")
            local version_type=${2:-"patch"}
            local next_version=$(get_next_version "$version_type")
            local message=${3:-"Release $next_version"}
            create_tag "$next_version" "$message"
            ;;
        "tag-custom")
            local custom_version=$2
            if [[ -z "$custom_version" ]]; then
                log_error "Custom version not specified"
                exit 1
            fi
            if [[ ! "$custom_version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+.*$ ]]; then
                log_error "Invalid version format. Use vX.Y.Z format"
                exit 1
            fi
            local message=${3:-"Release $custom_version"}
            create_tag "$custom_version" "$message"
            ;;
        "delete")
            delete_tag "$2"
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [command] [options]"
            echo ""
            echo "Commands:"
            echo "  info           - Show current version information (default)"
            echo "  current        - Show current version information"
            echo "  list           - List all tags"
            echo "  tag [type]     - Create a new tag (type: major, minor, patch)"
            echo "  tag-custom <v> - Create a custom version tag"
            echo "  delete <v>     - Delete a tag"
            echo "  help           - Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 info"
            echo "  $0 tag patch"
            echo "  $0 tag minor"
            echo "  $0 tag-custom v1.2.3"
            echo "  $0 delete v1.2.3"
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