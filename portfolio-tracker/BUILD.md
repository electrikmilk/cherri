# Build Instructions

This document explains how to build the Portfolio Tracker shortcuts from source.

## Prerequisites

1. **Go 1.19+** - Required to build the Cherri compiler
2. **Git** - To clone the repository
3. **Unix-like environment** - macOS, Linux, or WSL on Windows

## Step 1: Build the Cherri Compiler

The Portfolio Tracker is built using Cherri, which needs to be compiled from source first.

```bash
# If you haven't already, clone the cherri repository
# (Skip if you're already in the cherri repo)
git clone https://github.com/electrikmilk/cherri.git
cd cherri

# Build the cherri compiler
go build -o cherri .

# Verify the build
./cherri --version
```

This will create a `cherri` executable in the repository root.

## Step 2: Compile Portfolio Tracker Shortcuts

Once the Cherri compiler is built, you can compile the Portfolio Tracker shortcuts:

```bash
# Navigate to the portfolio-tracker directory
cd portfolio-tracker

# Run the build script
./build.sh
```

The build script will:
- Check that the cherri compiler is available
- Compile all 5 shortcuts from `src/shortcuts/`
- Output compiled `.shortcut` files to `dist/`

### Manual Compilation

If you prefer to compile shortcuts manually:

```bash
# Navigate to the cherri repository root
cd /path/to/cherri

# Compile each shortcut
./cherri portfolio-tracker/src/shortcuts/add_holding.cherri \
  -o portfolio-tracker/dist/Portfolio_Add_Holding.shortcut

./cherri portfolio-tracker/src/shortcuts/update_holding.cherri \
  -o portfolio-tracker/dist/Portfolio_Update_Holding.shortcut

./cherri portfolio-tracker/src/shortcuts/remove_holding.cherri \
  -o portfolio-tracker/dist/Portfolio_Remove_Holding.shortcut

./cherri portfolio-tracker/src/shortcuts/view_portfolio.cherri \
  -o portfolio-tracker/dist/Portfolio_View.shortcut

./cherri portfolio-tracker/src/shortcuts/refresh_widget.cherri \
  -o portfolio-tracker/dist/Portfolio_Refresh_Widget.shortcut
```

## Step 3: Transfer to iOS Device

After compilation, transfer the `.shortcut` files to your iOS device:

### Option 1: AirDrop (macOS)
1. Open Finder
2. Navigate to `portfolio-tracker/dist/`
3. Select all 5 `.shortcut` files
4. Right-click → Share → AirDrop
5. Select your iPhone/iPad

### Option 2: iCloud Drive
1. Upload all `.shortcut` files to iCloud Drive
2. On iOS: Open Files app → iCloud Drive
3. Tap each `.shortcut` file to import

### Option 3: Email/Messages
1. Attach `.shortcut` files to an email or message
2. Send to yourself
3. Open on iOS device and tap each attachment

## Step 4: Import into Shortcuts

On your iOS device:
1. Tap each `.shortcut` file
2. iOS will open the Shortcuts app
3. Tap "Add Shortcut" for each one
4. Grant any requested permissions

## Verification

After importing, you should have 5 new shortcuts:
- ✅ Portfolio - Add Holding
- ✅ Portfolio - Update Holding
- ✅ Portfolio - Remove Holding
- ✅ Portfolio - View
- ✅ Portfolio - Refresh Widget

Run "Portfolio - Add Holding" to test the installation.

## Troubleshooting

### "cherri: command not found"
- Make sure you built the cherri compiler (Step 1)
- Ensure you're running the build script from the portfolio-tracker directory
- The cherri binary must be in your PATH or in the parent directory

### Compilation errors
- Check that all library files exist in `src/lib/`
- Verify that all shortcuts reference the correct library paths
- Look for syntax errors in the `.cherri` files

### Build script fails
- Ensure the script is executable: `chmod +x build.sh`
- Check that Go is properly installed: `go version`
- Verify you have write permissions to the `dist/` directory

### Import fails on iOS
- Make sure the `.shortcut` files were compiled successfully
- Try importing one at a time
- Check iOS version compatibility (iOS 15+)

## Development

### Modifying Shortcuts

To modify the shortcuts:

1. Edit files in `src/shortcuts/` or `src/lib/`
2. Rebuild: `./build.sh`
3. Transfer updated `.shortcut` files to iOS
4. Re-import (they will update existing shortcuts)

### Testing

It's recommended to test on a real iOS device, as the Shortcuts app behavior can differ from simulators.

### Adding New Shortcuts

To add a new shortcut:

1. Create `src/shortcuts/your_shortcut.cherri`
2. Add appropriate metadata:
   ```cherri
   #define name "Your Shortcut Name"
   #define color blue
   #define glyph star
   ```
3. Update `build.sh` to include your new shortcut
4. Run `./build.sh`

## Next Steps

After building, see [SETUP.md](SETUP.md) for end-user setup instructions.

## Resources

- [Cherri Documentation](https://cherrilang.org)
- [Cherri GitHub Repository](https://github.com/electrikmilk/cherri)
- [Apple Shortcuts User Guide](https://support.apple.com/guide/shortcuts/)
