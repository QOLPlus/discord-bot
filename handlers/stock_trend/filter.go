package stock_trend

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
	case "GERMAN":
		return ":flag_de:"
	default:
		return ""
	}
}
