package cryptocurrency

import (
	coreCrypto "github.com/QOLPlus/core/commands/cryptocurrency"
)

func arrowOfTicker(ticker coreCrypto.Ticker) string {
	switch ticker.Change {
	case "EVEN":
		return "-"
	case "RISE":
		return "▲"
	case "FALL":
		return "▼"
	default:
		return ""
	}
}
func pmChangeRate(ticker coreCrypto.Ticker) float64 {
	switch ticker.Change {
	case "FALL":
		return -1 * ticker.ChangeRate
	default:
		return ticker.ChangeRate
	}
}
