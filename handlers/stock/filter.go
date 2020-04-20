package stock

import (
	"github.com/QOLPlus/core/commands/stock"
)

func arrowOfSecurity(security stock.FetchedSecurity) string {
	switch security.Change {
	case "UPPER_LIMIT":
		return "↑"
	case "RISE":
		return "▲"
	case "FALL":
		return "▼"
	case "LOWER_LIMIT":
		return "↓"
	default:
		return ""
	}
}

func colorOfSecurity(security stock.FetchedSecurity) int {
	switch security.Change {
	case "UPPER_LIMIT", "RISE":
		return 16062488
	case "FALL", "LOWER_LIMIT":
		return 1794513
	default:
		return 0
	}
}

func asterNameIfKosdaq(security stock.FetchedSecurity) string {
	aster := ""
	if security.Market == "KOSDAQ" {
		aster = "*"
	}
	return security.Name + aster
}
