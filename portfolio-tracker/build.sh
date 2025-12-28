#!/bin/bash

# Portfolio Tracker - Build Script
# Compiles all Cherri shortcuts to .shortcut files

set -e

echo "Portfolio Tracker - Build Script"
echo "================================="
echo ""

# Check if cherri is available
if ! command -v cherri &> /dev/null; then
    echo "Error: cherri compiler not found"
    echo "Please build cherri first by running: go build -o cherri ."
    echo "Make sure you're in the cherri repository root directory"
    exit 1
fi

# Navigate to project directory
cd "$(dirname "$0")"
PROJECT_DIR=$(pwd)
SRC_DIR="$PROJECT_DIR/src"
DIST_DIR="$PROJECT_DIR/dist"

# Create dist directory if it doesn't exist
mkdir -p "$DIST_DIR"

echo "Source directory: $SRC_DIR"
echo "Output directory: $DIST_DIR"
echo ""

# Array of shortcuts to compile
declare -a SHORTCUTS=(
    "add_holding:Portfolio_Add_Holding"
    "update_holding:Portfolio_Update_Holding"
    "remove_holding:Portfolio_Remove_Holding"
    "view_portfolio:Portfolio_View"
    "refresh_widget:Portfolio_Refresh_Widget"
)

# Compile each shortcut
for item in "${SHORTCUTS[@]}"; do
    IFS=':' read -r source output <<< "$item"

    SOURCE_FILE="$SRC_DIR/shortcuts/${source}.cherri"
    OUTPUT_FILE="$DIST_DIR/${output}.shortcut"

    echo "Compiling: ${source}.cherri -> ${output}.shortcut"

    if [ ! -f "$SOURCE_FILE" ]; then
        echo "  Error: Source file not found: $SOURCE_FILE"
        exit 1
    fi

    cherri "$SOURCE_FILE" -o "$OUTPUT_FILE"

    if [ $? -eq 0 ]; then
        echo "  ✓ Success"
    else
        echo "  ✗ Failed"
        exit 1
    fi
done

echo ""
echo "================================="
echo "Build complete!"
echo ""
echo "Generated files:"
ls -lh "$DIST_DIR"/*.shortcut
echo ""
echo "To install:"
echo "1. Transfer the .shortcut files to your iOS device"
echo "2. Open each file to import into Shortcuts app"
echo "3. See SETUP.md for complete setup instructions"
