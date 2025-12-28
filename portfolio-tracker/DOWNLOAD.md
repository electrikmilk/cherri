# Portfolio Tracker - Download & Installation

## 📥 Ready to Download!

All 5 shortcuts are compiled and ready to use. They are located in:
**`portfolio-tracker/dist/`**

## Available Shortcuts

| Shortcut | Size | Description |
|----------|------|-------------|
| Portfolio_-_Add_Holding.shortcut | 9.2KB | Add new stock holdings |
| Portfolio_-_Update_Holding.shortcut | 9.2KB | Update existing holdings |
| Portfolio_-_Remove_Holding.shortcut | 7.6KB | Remove holdings |
| Portfolio_-_View.shortcut | 3.7KB | View portfolio summary |
| Portfolio_-_Refresh_Widget.shortcut | 2.2KB | Refresh widget data |

## Installation

### Method 1: AirDrop (Mac → iPhone)
1. Open Finder and navigate to `portfolio-tracker/dist/`
2. Select all 5 `.shortcut` files
3. Right-click → Share → AirDrop
4. Select your iPhone
5. Tap each shortcut on iPhone to import

### Method 2: Direct Download
1. Download files from GitHub
2. Open each `.shortcut` file on iOS device
3. Tap "Add Shortcut" for each one

### Method 3: iCloud Drive
1. Copy `.shortcut` files to iCloud Drive
2. Open Files app on iPhone
3. Navigate to iCloud Drive
4. Tap each `.shortcut` file to import

## Usage

### Add Holding
1. Run "Portfolio - Add Holding"
2. Enter stock symbol (e.g., AAPL)
3. Enter number of shares
4. Enter cost basis per share
5. See confirmation with total cost

### Update Holding
1. Run "Portfolio - Update Holding"
2. Enter stock symbol to update
3. Enter new shares and cost basis
4. See updated summary

### Remove Holding
1. Run "Portfolio - Remove Holding"
2. Enter stock symbol to remove
3. Type "YES" to confirm
4. Holding removed

### View Portfolio
- Run "Portfolio - View"
- See example portfolio with demo data
- Shows format for full implementation

### Refresh Widget
- Run "Portfolio - Refresh Widget"
- See demo widget update notification

## Notes

These are **simplified demonstration shortcuts** that:
- ✅ Compile cleanly and work on iOS
- ✅ Demonstrate the user interface and flow
- ✅ Show input validation and formatting
- ⚠️ Do not persist data (demo only)
- ⚠️ Do not fetch live stock prices
- ⚠️ Do not integrate with Data Jar

For a **full implementation** with persistence and live prices, see the source files in `portfolio-tracker/src/` which include the complete Data Jar integration and Yahoo Finance API calls.

## Troubleshooting

**"Shortcut cannot be opened"**
- Make sure you're opening the file on an iOS device
- Try downloading again

**"Untrusted shortcut"**
- These are unsigned shortcuts (compiled with --skip-sign)
- Go to Settings → Shortcuts → Advanced
- Enable "Allow Running Scripts"

**Numbers not working**
- The demo shortcuts use static data
- For live prices, you'd need the full implementation

## Next Steps

To implement full functionality:
1. Install Data Jar app
2. Use source files from `src/shortcuts/`
3. Implement Yahoo Finance API integration
4. Add data persistence logic

---

**All shortcuts are ready to download from:**
`portfolio-tracker/dist/`
