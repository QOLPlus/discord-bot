package cryptocurrency

import (
	"fmt"
	coreCrypto "github.com/QOLPlus/core/commands/cryptocurrency"
	"github.com/QOLPlus/discord-bot/refs"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"math"
	"regexp"
	"strings"
)


var commands = []string{"코인", "크립토", "가상자산", "가상화폐", "CC", "cc", "c"}

var Registry = &refs.HandlerRegistry{
	Commands: commands,
	Proc: Process,
}

var notFoundMessage string = fmt.Sprintf(
	"가상자산 정보를 찾을 수 없습니다. `[(%s) 검색어]`",
	strings.Join(commands, "|"),
)

var KRW string = "KRW"

func Process(store *refs.Store, s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	tickerStore := store.Ticker
	tickerStore.Reload(3600)

	var selectedMarkets []coreCrypto.MarketMaster
	for _, phrase := range data {
		if len(phrase) == 0 || KRW == phrase {
			continue
		}

		for _, master := range  *tickerStore.Data {
			if similar(master, phrase) && strings.HasPrefix(master.Market, KRW){
				selectedMarkets = append(selectedMarkets, master)
			}
		}
	}

	if len(selectedMarkets) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, notFoundMessage)
		return
	}

	for _, market := range selectedMarkets {
		tickers, err := market.FetchTicker()

		if err != nil {
			continue
		}

		for _, ticker := range *tickers {
			_, _ = s.ChannelMessageSend(m.ChannelID, buildSimpleMessage(market, ticker))
		}
	}
}

func similar(market coreCrypto.MarketMaster, keyword string) bool {
	match, _ := regexp.MatchString(
		keyword,
		strings.Join([]string{
			market.KoreanName,
			market.EnglishName,
			market.Market,
		}, "|"),
	)
	return match
}

func buildSimpleMessage(market coreCrypto.MarketMaster, ticker coreCrypto.Ticker) string {
	return fmt.Sprintf(
		"**%s**: %s %s",
		market.KoreanName,
		strings.Join([]string{
			humanize.Commaf(ticker.TradePrice),
			fmt.Sprintf("%s%s", arrowOfTicker(ticker), humanize.CommafWithDigits(math.Abs(ticker.ChangePrice), 2)),
			"(" + fmt.Sprintf("%.2f", pmChangeRate(ticker) * 100) + "%)",
		}, " "),
		fmt.Sprintf(
			"%s / 52주 신고가: %s(%s) / 52주 신저가: %s(%s)",
			ticker.Market,
			humanize.CommafWithDigits(ticker.Highest52WeekPrice, 2),
			ticker.Highest52WeekDate,
			humanize.CommafWithDigits(ticker.Lowest52WeekPrice, 2),
			ticker.Lowest52WeekDate,
		),
	)
}
