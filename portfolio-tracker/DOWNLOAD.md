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
| Portfolio_-_View.shortcut | 49KB | **View portfolio with LIVE prices** 📈 |
| Portfolio_-_Refresh_Widget.shortcut | 28KB | **Refresh with LIVE prices** 📈 |

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
- **Fetches LIVE stock prices from Yahoo Finance API** 📈
- Shows example portfolio with real-time AAPL and GOOGL prices
- Calculates actual gains/losses based on current market prices

### Refresh Widget
- Run "Portfolio - Refresh Widget"
- **Fetches LIVE prices for portfolio stocks** 📈
- Shows current portfolio value with real-time data
- Displays actual market prices for AAPL and GOOGL

## Features

These working shortcuts include:
- ✅ Compile cleanly and work on iOS
- ✅ User-friendly interface with input validation
- ✅ **LIVE stock price fetching via Yahoo Finance API** 📈
- ✅ Real-time portfolio value calculations
- ✅ Actual gain/loss calculations with live data
- ⚠️ Do not persist data (demo portfolio only)
- ⚠️ Do not integrate with Data Jar (uses example holdings)

**Example Portfolio (hardcoded for demo):**
- AAPL: 100 shares @ $150.25 cost basis
- GOOGL: 50 shares @ $135.00 cost basis
- Fetches current market prices and calculates real returns!

For a **full implementation** with data persistence, see the source files in `portfolio-tracker/src/` which include complete Data Jar integration.

## Troubleshooting

**"Shortcut cannot be opened"**
- Make sure you're opening the file on an iOS device
- Try downloading again

**"Untrusted shortcut"**
- These are unsigned shortcuts (compiled with --skip-sign)
- Go to Settings → Shortcuts → Advanced
- Enable "Allow Running Scripts"
- Enable "Allow Sharing Large Amounts of Data"

**"Cannot connect to Yahoo Finance"**
- Make sure you have internet connection
- View Portfolio and Refresh Widget require network access
- If API is down, shortcuts will show error or 0 prices
- Yahoo Finance API is free and usually very reliable

**"Prices seem wrong"**
- Prices are LIVE from Yahoo Finance API
- They reflect actual current market prices
- Market must be open for real-time updates
- After hours, shows last closing price

## Next Steps

To implement full persistence functionality:
1. Install Data Jar app
2. See source files in `src/shortcuts/` for reference
3. **Live price fetching is already implemented!** ✅
4. Add data persistence logic to store/retrieve holdings

---

**All shortcuts are ready to download from:**
`portfolio-tracker/dist/`
