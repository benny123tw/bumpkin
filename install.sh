#!/bin/sh
set -e

# bumpkin installer script
# Usage: curl -sSfL https://raw.githubusercontent.com/benny123tw/bumpkin/main/install.sh | sh -s -- -b /usr/local/bin

usage() {
  cat <<EOF
$0: download bumpkin binary

Usage: $0 [-b <bindir>] [-d] [<tag>]
  -b  installation directory, defaults to ./bin
  -d  enable debug logging
  <tag>  version tag (e.g., v0.5.1), defaults to latest

Examples:
  $0                           # install latest to ./bin
  $0 -b /usr/local/bin         # install latest to /usr/local/bin
  $0 -b \$HOME/.local/bin v0.5.0  # install specific version

EOF
  exit 2
}

log_info() { echo "[info] $*" >&2; }
log_err() { echo "[error] $*" >&2; }
log_debug() { [ "$DEBUG" = "1" ] && echo "[debug] $*" >&2 || true; }

# Parse arguments
BINDIR=${BINDIR:-./bin}
DEBUG=0
while getopts "b:dh?" arg; do
  case "$arg" in
    b) BINDIR="$OPTARG" ;;
    d) DEBUG=1 ;;
    h | \?) usage ;;
  esac
done
shift $((OPTIND - 1))
TAG=$1

# Detect OS and architecture
detect_platform() {
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m)

  case "$OS" in
    msys* | mingw* | cygwin*) OS="windows" ;;
  esac

  case "$ARCH" in
    x86_64 | amd64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
    *)
      log_err "unsupported architecture: $ARCH"
      exit 1
      ;;
  esac

  # windows/arm64 is not supported
  if [ "$OS" = "windows" ] && [ "$ARCH" = "arm64" ]; then
    log_err "windows/arm64 is not supported"
    exit 1
  fi

  log_debug "detected platform: $OS/$ARCH"
}

# Get latest release tag from GitHub
get_latest_tag() {
  if command -v curl >/dev/null 2>&1; then
    curl -sSfL "https://api.github.com/repos/benny123tw/bumpkin/releases/latest" |
      grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "https://api.github.com/repos/benny123tw/bumpkin/releases/latest" |
      grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
  else
    log_err "curl or wget required"
    exit 1
  fi
}

# Download file
download() {
  url=$1
  dest=$2
  log_debug "downloading $url"
  if command -v curl >/dev/null 2>&1; then
    curl -sSfL -o "$dest" "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -q -O "$dest" "$url"
  else
    log_err "curl or wget required"
    exit 1
  fi
}

# Verify checksum
verify_checksum() {
  file=$1
  checksums=$2
  filename=$(basename "$file")

  want=$(grep "$filename" "$checksums" | cut -d ' ' -f 1)
  if [ -z "$want" ]; then
    log_err "checksum not found for $filename"
    return 1
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    got=$(sha256sum "$file" | cut -d ' ' -f 1)
  elif command -v shasum >/dev/null 2>&1; then
    got=$(shasum -a 256 "$file" | cut -d ' ' -f 1)
  else
    log_err "sha256sum or shasum required for verification"
    return 1
  fi

  if [ "$want" != "$got" ]; then
    log_err "checksum mismatch: expected $want, got $got"
    return 1
  fi
  log_debug "checksum verified"
}

# Main installation
main() {
  detect_platform

  # Get version
  if [ -z "$TAG" ]; then
    log_info "fetching latest version..."
    TAG=$(get_latest_tag)
    if [ -z "$TAG" ]; then
      log_err "failed to get latest version"
      exit 1
    fi
  fi
  VERSION=${TAG#v}
  log_info "installing bumpkin $TAG for $OS/$ARCH"

  # Construct download URLs
  # Binary format: bumpkin-{version}-{os}-{arch}
  BINARY_NAME="bumpkin-${VERSION}-${OS}-${ARCH}"
  [ "$OS" = "windows" ] && BINARY_NAME="${BINARY_NAME}.exe"

  BASE_URL="https://github.com/benny123tw/bumpkin/releases/download/${TAG}"
  BINARY_URL="${BASE_URL}/${BINARY_NAME}"
  CHECKSUM_URL="${BASE_URL}/checksums.txt"

  # Create temp directory
  tmpdir=$(mktemp -d)
  trap 'rm -rf "$tmpdir"' EXIT

  # Download binary and checksums
  log_info "downloading binary..."
  download "$BINARY_URL" "$tmpdir/$BINARY_NAME"
  download "$CHECKSUM_URL" "$tmpdir/checksums.txt"

  # Verify checksum
  log_info "verifying checksum..."
  verify_checksum "$tmpdir/$BINARY_NAME" "$tmpdir/checksums.txt"

  # Install binary
  mkdir -p "$BINDIR"
  DEST_NAME="bumpkin"
  [ "$OS" = "windows" ] && DEST_NAME="bumpkin.exe"

  install -m 755 "$tmpdir/$BINARY_NAME" "$BINDIR/$DEST_NAME"
  log_info "installed $BINDIR/$DEST_NAME"

  # Verify installation
  if [ -x "$BINDIR/$DEST_NAME" ]; then
    log_info "installation complete!"
    "$BINDIR/$DEST_NAME" version 2>/dev/null || true
  fi
}

main
