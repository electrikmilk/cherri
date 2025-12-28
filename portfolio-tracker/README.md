# Portfolio Tracker - Cherri Implementation

A comprehensive stock portfolio tracking system built with **Cherri** (Siri Shortcuts programming language). Track your investments, monitor gains/losses, and display real-time portfolio data on your iOS/macOS home screen.

## Features

- **Add Holdings**: Add stocks with symbol, shares, and cost basis
- **Update Holdings**: Modify existing positions
- **Remove Holdings**: Delete holdings from your portfolio
- **View Portfolio**: See complete portfolio summary with current prices and gains/losses
- **Widget Support**: Real-time portfolio widget for iOS/macOS home screen
- **Auto-Refresh**: Optional daily automation to keep data current
- **Yahoo Finance Integration**: Live stock price data
- **Data Jar Storage**: Persistent, iCloud-synced data storage

## Components

### Shortcuts (5)
1. **Portfolio - Add Holding** - Add new stock positions
2. **Portfolio - Update Holding** - Modify existing positions
3. **Portfolio - Remove Holding** - Delete positions
4. **Portfolio - View** - Display full portfolio summary
5. **Portfolio - Refresh Widget** - Update widget data cache

### Widget (1)
- **Portfolio Widget** (Scriptable) - Home screen widget displaying portfolio value, change, and top holdings

## Quick Start

See [SETUP.md](SETUP.md) for complete installation and setup instructions.

### Prerequisites
- iOS/iPadOS 15+ or macOS Sonoma+
- Data Jar (free app)
- Scriptable (free app)
- iCloud Drive enabled

### Installation Summary
1. Install Data Jar and Scriptable from App Store
2. Download and install the 5 compiled `.shortcut` files
3. Copy `Portfolio_Widget.js` to Scriptable
4. Add Scriptable widget to home screen
5. Run "Portfolio - Add Holding" to add your first stock

## Usage

### Adding a Stock
1. Run "Portfolio - Add Holding"
2. Enter stock symbol (e.g., AAPL)
3. Enter number of shares
4. Enter cost basis per share

### Viewing Your Portfolio
1. Run "Portfolio - View" to see full details
2. Check your home screen widget for quick overview

### Updating the Widget
- Widget updates automatically when you add/update/remove holdings
- Run "Portfolio - Refresh Widget" manually for latest prices
- Set up automation for daily auto-refresh

## Architecture

### Data Storage
Uses **Data Jar** for persistent storage with the following structure:

```json
{
  "holdings": [
    {
      "symbol": "AAPL",
      "shares": 100,
      "costBasis": 150.25,
      "dateAdded": "2024-01-15"
    }
  ],
  "lastUpdated": "2024-12-28T09:00:00Z"
}
```

### Yahoo Finance API
Fetches real-time stock prices from:
```
https://query1.finance.yahoo.com/v8/finance/chart/{SYMBOL}
```

## Project Structure

```
portfolio-tracker/
├── src/
│   ├── lib/
│   │   ├── portfolio_data.cherri    # Data access layer
│   │   ├── yahoo_api.cherri         # API integration
│   │   └── formatting.cherri        # Display utilities
│   └── shortcuts/
│       ├── add_holding.cherri
│       ├── update_holding.cherri
│       ├── remove_holding.cherri
│       ├── view_portfolio.cherri
│       └── refresh_widget.cherri
├── widget/
│   └── Portfolio_Widget.js          # Scriptable widget
├── dist/                             # Compiled shortcuts
├── SETUP.md                          # Setup instructions
└── README.md                         # This file
```

## Building from Source

### Requirements
- Cherri compiler (https://cherrilang.org)
- Git

### Compilation

```bash
cd portfolio-tracker/src

# Compile all shortcuts
cherri shortcuts/add_holding.cherri -o ../dist/Portfolio_Add_Holding.shortcut
cherri shortcuts/update_holding.cherri -o ../dist/Portfolio_Update_Holding.shortcut
cherri shortcuts/remove_holding.cherri -o ../dist/Portfolio_Remove_Holding.shortcut
cherri shortcuts/view_portfolio.cherri -o ../dist/Portfolio_View.shortcut
cherri shortcuts/refresh_widget.cherri -o ../dist/Portfolio_Refresh_Widget.shortcut
```

## Troubleshooting

### Widget Not Updating
- Ensure internet connection is active
- Run "Portfolio - Refresh Widget" manually
- Check that Data Jar is installed and accessible

### Price Data Issues
- Yahoo Finance API may have temporary outages
- Some symbols may not be supported (try alternative exchanges)
- Rate limiting may occur with frequent requests

### Data Jar Connection
- Open Data Jar at least once to initialize
- Grant necessary permissions in iOS Settings
- Ensure iCloud is enabled

## Limitations

- Yahoo Finance API is free but may have rate limits
- Widget refresh requires manual trigger or automation
- International stocks may require exchange suffix (e.g., AAPL.L)
- Historical data not stored (only current prices)

## Future Enhancements

- [ ] Support for cryptocurrency portfolios
- [ ] Historical performance charts
- [ ] Dividend tracking
- [ ] Multiple portfolio support
- [ ] Export to CSV
- [ ] Tax lot management

## License

MIT License - See LICENSE file for details

## Credits

- Built with [Cherri](https://cherrilang.org)
- Stock data from [Yahoo Finance](https://finance.yahoo.com)
- Uses [Data Jar](https://datajar.app) for storage
- Widget powered by [Scriptable](https://scriptable.app)

## Support

For issues, questions, or contributions:
- Open an issue on GitHub
- Check SETUP.md for common problems
- Review Cherri documentation at https://cherrilang.org

## Version

Current Version: **1.0.0**

---

Made with ❤️ using Cherri
