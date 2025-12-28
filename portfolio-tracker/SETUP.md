# Portfolio Tracker Setup Guide

## Prerequisites
1. **Data Jar** (free) - App Store: https://apps.apple.com/app/data-jar/id1453273600
2. **Scriptable** (free) - App Store: https://apps.apple.com/app/scriptable/id1405459188
3. **iCloud Drive** enabled on your device

## Step 1: Install Shortcuts

1. Download all 5 `.shortcut` files from the releases
2. Open each file - it will open in the Shortcuts app
3. Tap "Add Shortcut" for each

The shortcuts are:
- Portfolio - Add Holding
- Portfolio - Update Holding
- Portfolio - Remove Holding
- Portfolio - View
- Portfolio - Refresh Widget

## Step 2: Install Scriptable Widget

1. Open the Scriptable app
2. Tap the "+" button to create a new script
3. Copy the entire contents of `Portfolio_Widget.js`
4. Paste into Scriptable
5. Rename the script to "Portfolio Widget"
6. Tap "Done"

## Step 3: Add Widget to Home Screen

### iOS:
1. Long-press on your home screen
2. Tap the "+" button (top left)
3. Search for "Scriptable"
4. Choose your preferred widget size (Medium recommended)
5. Tap "Add Widget"
6. Long-press the new widget → "Edit Widget"
7. Set Script to "Portfolio Widget"
8. Set "When Interacting" to "Run Script"

### macOS (Sonoma+):
1. Complete iOS setup first
2. On Mac: Right-click desktop → "Edit Widgets"
3. Find Scriptable (marked "From iPhone")
4. Drag to desktop

## Step 4: Add Your First Holding

1. Run "Portfolio - Add Holding" (via Shortcuts app, Siri, or widget)
2. Enter: Stock symbol (e.g., AAPL)
3. Enter: Number of shares
4. Enter: Cost basis per share
5. Done! Your widget will update automatically.

## Step 5: Setup Daily Auto-Refresh (Optional)

1. Open Shortcuts app
2. Go to "Automation" tab
3. Tap "+" → "Personal Automation"
4. Select "Time of Day"
5. Set to 9:00 AM, Daily
6. Tap "Next"
7. Add Action → "Run Shortcut" → "Portfolio - Refresh Widget"
8. Tap "Next"
9. Toggle OFF "Ask Before Running"
10. Tap "Done"

## Troubleshooting

**Widget shows "No holdings":**
- Run "Portfolio - Add Holding" first
- Make sure Data Jar is installed

**Prices not updating:**
- Check internet connection
- Yahoo Finance API may be temporarily unavailable
- Try running "Portfolio - Refresh Widget" manually

**Widget not appearing on macOS:**
- Requires macOS Sonoma or later
- iPhone must be on same iCloud account and nearby

**Data Jar errors:**
- Ensure Data Jar is installed and granted necessary permissions
- Open Data Jar once to initialize it

**Compilation errors:**
- Make sure you have the latest version of Cherri installed
- Check that all library files are in the correct locations
