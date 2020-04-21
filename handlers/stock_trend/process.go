package stock_trend

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"

	"github.com/QOLPlus/core/commands/stock"
	"github.com/QOLPlus/discord-bot/refs"
)


var Registry = &refs.HandlerRegistry{
	Commands: []string{"시황"},
	Proc:     Process,
}

type exchangeMap map[string][]stock.FetchedSecurity
func (e exchangeMap) buildMessage() string {
	var globalMessages []string
	for country, securities := range e {
		var countryMessages []string
		for _, security := range securities {
			countryMessages = append(
				countryMessages,
				fmt.Sprintf(
					"%s **%s**: %.2f %s",
					flagOfCountry(country),
					security.Name,
					security.TradePrice,
					"(" + fmt.Sprintf("%.2f", security.ChangePriceRate * 100) + "%)",
				),
			)
		}

		globalMessages = append(
			globalMessages,
			strings.Join(countryMessages, "  "),
		)
	}
	return strings.Join(globalMessages, "\n")
}

func Process(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	fetched, err := stock.FetchTrends(nil)
	if err != nil {
		return
	}

	data := exchangeMap{}
	for _, security := range fetched.RecentSecurities {
		data[security.ExchangeCountry] = append(data[security.ExchangeCountry], security)
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, data.buildMessage())
}

