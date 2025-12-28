// Portfolio Widget for Scriptable
// Reads data from Data Jar (set by Shortcuts)

const WIDGET_VERSION = "1.0.0";

// Data Jar integration
async function getDataJarValue(key) {
  const url = `datajar://x-callback-url/get?key=${encodeURIComponent(key)}`;

  // For Scriptable, we read from shared file or use Keychain as fallback
  // Data Jar doesn't have direct Scriptable integration, so we use shared file

  try {
    const fm = FileManager.iCloud();
    const path = fm.joinPath(fm.documentsDirectory(), "portfolio_widget_data.json");

    if (fm.fileExists(path)) {
      await fm.downloadFileFromiCloud(path);
      const data = fm.readString(path);
      return JSON.parse(data);
    }
  } catch (e) {
    console.log("Error reading data: " + e);
  }

  return null;
}

// Alternative: Read directly from portfolio.csv if using CSV storage
async function readPortfolioCSV() {
  const fm = FileManager.iCloud();
  const shortcutsDir = fm.joinPath(fm.documentsDirectory(), "../Shortcuts");
  const csvPath = fm.joinPath(shortcutsDir, "portfolio.csv");

  if (!fm.fileExists(csvPath)) {
    return [];
  }

  await fm.downloadFileFromiCloud(csvPath);
  const csvContent = fm.readString(csvPath);
  const lines = csvContent.split("\n").filter(l => l.trim());

  // Skip header
  const holdings = [];
  for (let i = 1; i < lines.length; i++) {
    const [symbol, shares, costBasis, dateAdded] = lines[i].split(",");
    if (symbol) {
      holdings.push({
        symbol: symbol.trim(),
        shares: parseFloat(shares),
        costBasis: parseFloat(costBasis),
        dateAdded: dateAdded?.trim() || ""
      });
    }
  }

  return holdings;
}

// Yahoo Finance API
async function fetchPrice(symbol) {
  const url = `https://query1.finance.yahoo.com/v8/finance/chart/${symbol}`;

  try {
    const req = new Request(url);
    const response = await req.loadJSON();

    if (response.chart?.result?.[0]?.meta) {
      return {
        price: response.chart.result[0].meta.regularMarketPrice,
        previousClose: response.chart.result[0].meta.previousClose
      };
    }
  } catch (e) {
    console.log(`Error fetching ${symbol}: ${e}`);
  }

  return { price: 0, previousClose: 0 };
}

async function fetchAllPrices(symbols) {
  const prices = {};

  for (const symbol of symbols) {
    prices[symbol] = await fetchPrice(symbol);
    // Small delay to avoid rate limiting
    await new Promise(r => setTimeout(r, 100));
  }

  return prices;
}

function calculatePortfolio(holdings, prices) {
  let totalValue = 0;
  let totalCost = 0;

  const enrichedHoldings = holdings.map(h => {
    const priceData = prices[h.symbol] || { price: 0, previousClose: 0 };
    const currentValue = h.shares * priceData.price;
    const costTotal = h.shares * h.costBasis;

    totalValue += currentValue;
    totalCost += costTotal;

    return {
      ...h,
      currentPrice: priceData.price,
      previousClose: priceData.previousClose,
      currentValue,
      costTotal,
      gainLoss: currentValue - costTotal,
      gainLossPercent: costTotal > 0 ? ((currentValue - costTotal) / costTotal) * 100 : 0
    };
  });

  return {
    holdings: enrichedHoldings,
    totalValue,
    totalCost,
    change: totalValue - totalCost,
    changePercent: totalCost > 0 ? ((totalValue - totalCost) / totalCost) * 100 : 0
  };
}

async function createWidget() {
  const widget = new ListWidget();
  widget.backgroundColor = new Color("#1a1a2e");
  widget.setPadding(12, 12, 12, 12);

  // Read holdings
  const holdings = await readPortfolioCSV();

  if (holdings.length === 0) {
    const text = widget.addText("📊 No holdings");
    text.textColor = Color.white();
    text.font = Font.mediumSystemFont(16);
    return widget;
  }

  const symbols = holdings.map(h => h.symbol);
  const prices = await fetchAllPrices(symbols);
  const portfolio = calculatePortfolio(holdings, prices);

  // Title
  const title = widget.addText("📊 PORTFOLIO");
  title.textColor = Color.white();
  title.font = Font.boldSystemFont(12);

  widget.addSpacer(4);

  // Total value
  const valueText = widget.addText(`$${portfolio.totalValue.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2})}`);
  valueText.textColor = Color.white();
  valueText.font = Font.boldSystemFont(24);

  // Change
  const changeColor = portfolio.change >= 0 ? Color.green() : Color.red();
  const changeSign = portfolio.change >= 0 ? "+" : "";
  const changeText = widget.addText(
    `${changeSign}$${Math.abs(portfolio.change).toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2})} (${changeSign}${portfolio.changePercent.toFixed(2)}%)`
  );
  changeText.textColor = changeColor;
  changeText.font = Font.mediumSystemFont(14);

  widget.addSpacer(8);

  // Holdings (for medium/large widgets)
  if (config.widgetFamily !== "small") {
    const displayHoldings = portfolio.holdings
      .sort((a, b) => b.currentValue - a.currentValue)
      .slice(0, 3);

    for (const holding of displayHoldings) {
      const line = widget.addText(
        `${holding.symbol}  ${holding.shares} @ $${holding.currentPrice.toFixed(2)}`
      );
      line.textColor = Color.white();
      line.font = Font.systemFont(12);
      line.lineLimit = 1;
    }
  }

  widget.addSpacer();

  // Last updated
  const now = new Date();
  const timeStr = now.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit'
  });
  const updated = widget.addText(`Updated: ${timeStr}`);
  updated.textColor = Color.gray();
  updated.font = Font.systemFont(10);

  // Tap to refresh
  widget.url = "shortcuts://run-shortcut?name=Portfolio%20-%20Refresh%20Widget";

  return widget;
}

// Run
const widget = await createWidget();
if (config.runsInWidget) {
  Script.setWidget(widget);
} else {
  widget.presentMedium();
}
Script.complete();
