package asset

import (
	coreAsset "github.com/QOLPlus/core/commands/asset"
)

func arrowOfSecurity(security coreAsset.Security) string {
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

func pmChangeRate(security coreAsset.Security) float64 {
	switch security.Change {
	case "FALL", "LOWER_LIMIT":
		return -1 * security.ChangeRate
	default:
		return security.ChangeRate
	}
}

func asterNameIfKosdaq(security coreAsset.Security) string {
	aster := ""
	if security.Market == "KOSDAQ" {
		aster = "*"
	}
	return security.Name + aster
}

func flagOfCountry(country string) string {
	switch country {
	case "KOREA":
		return ":flag_kr:"
	case "USA":
		return ":flag_us:"
	case "JAPAN":
		return ":flag_jp:"
	case "BRITISH":
		return ":flag_gb:"
	case "SHANGHAI":
		return ":flag_cn:"
	case "HONGKONG":
		return ":flag_hk:"
	case "GERMAN":
		return ":flag_de:"
	case "INDIA":
		return ":flag_in:"
	case "FRANCE":
		return ":flag_fr:"
	case "TAIWAN":
		return ":flag_tw:"
	case "ARGENTINA":
		return ":flag_ar:"
	case "BRAZIL":
		return ":flag_br:"
	case "HUNGARY":
		return ":flag_hu:"
	case "INDONESIA":
		return ":flag_id:"
	case "ITALY":
		return ":flag_it:"
	case "VIETNAM":
		return ":flag_vn:"
	case "MALAYSIA":
		return ":flag_my:"
	case "SINGAPORE":
		return ":flag_sg:"
	case "RUSIA":
		return ":flag_ru:"
	default:
		return "#" + country
	}
}
