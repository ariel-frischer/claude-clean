#!/bin/bash
set -e

# Build binaries for multiple platforms
# Usage: ./scripts/build-binaries.sh [version]

VERSION=${1:-"dev"}
BINARY_NAME="cclean"
OUTPUT_DIR="release"

echo "Building ${BINARY_NAME} binaries (version: ${VERSION})..."

# Clean and create output directory
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"

# Build for each platform
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r -a parts <<< "$platform"
    GOOS="${parts[0]}"
    GOARCH="${parts[1]}"

    output_name="${BINARY_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi

    echo "  Building ${output_name}..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "${OUTPUT_DIR}/${output_name}" .
done

echo ""
echo "✓ Binaries built successfully in ${OUTPUT_DIR}/"
ls -lh "${OUTPUT_DIR}"

# Generate checksums (cross-platform: Linux uses sha256sum, macOS uses shasum)
echo ""
echo "Generating checksums..."
cd "${OUTPUT_DIR}"
if command -v sha256sum &> /dev/null; then
    sha256sum ${BINARY_NAME}-* > SHA256SUMS
else
    shasum -a 256 ${BINARY_NAME}-* > SHA256SUMS
fi
cd ..

echo "✓ Checksums generated: ${OUTPUT_DIR}/SHA256SUMS"
echo ""
echo "Done! Binaries ready for release."
