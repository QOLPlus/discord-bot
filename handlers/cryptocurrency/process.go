package cryptocurrency

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"math"
	"regexp"
	"strings"

	"github.com/QOLPlus/core/commands/cryptocurrency"
	"github.com/QOLPlus/discord-bot/refs"
)

var commands = []string{"코인", "크립토", "가상자산", "가상화폐", "CC", "cc", "c"}

var Registry = &refs.HandlerRegistry{
	Commands: commands,
	Proc:     Process,
}

var notFoundMessage string = fmt.Sprintf(
	"가상자산 정보를 찾을 수 없습니다. `[(%s) 검색어]`",
	strings.Join(commands, "|"),
)

const (
	KRW      string = "KRW"
	MaxCount int    = 3
)

func Process(store *refs.Store, s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	tickerStore := store.Ticker
	tickerStore.Reload(3600)

	var markets []string
	for _, phrase := range data {
		if len(phrase) == 0 || KRW == phrase {
			continue
		}

		var selectedMarkets []cryptocurrency.MarketMaster
		for _, master := range *tickerStore.Data {
			if similar(master, phrase) && strings.HasPrefix(master.Market, KRW) {
				selectedMarkets = append(selectedMarkets, master)
			}
		}
		if len(selectedMarkets) == 0 {
			continue
		}

		var tMarkets []string
		for _, market := range selectedMarkets {
			if include(market, phrase) {
				tMarkets = append(tMarkets, market.Market)
				break
			}
		}

		if len(tMarkets) == 0 {
			tMarkets = append(tMarkets, selectedMarkets[0].Market)
		}

		markets = append(markets, tMarkets...)
	}

	if len(markets) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, notFoundMessage)
		return
	}

	tickers, err := cryptocurrency.FetchTickerByMarkets(markets)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, notFoundMessage)
		return
	}

	for idx, ticker := range *tickers {
		if idx == MaxCount {
			warning := fmt.Sprintf("한번에 %d종목 이상을 조회할 수 없습니다.", MaxCount)
			_, _ = s.ChannelMessageSend(m.ChannelID, warning)
			break
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, buildSimpleMessage(ticker))
	}
}

func similar(market cryptocurrency.MarketMaster, keyword string) bool {
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

func include(market cryptocurrency.MarketMaster, keyword string) bool {
	for _, token := range []string{
		market.KoreanName,
		market.EnglishName,
		market.Market,
	} {
		if token == keyword {
			return true
		}
	}

	return false
}

func buildSimpleMessage(ticker cryptocurrency.Ticker) string {
	return fmt.Sprintf(
		"**%s**: %s %s",
		ticker.Market,
		strings.Join([]string{
			humanize.Commaf(ticker.TradePrice),
			fmt.Sprintf("%s%s", arrowOfTicker(ticker), humanize.CommafWithDigits(math.Abs(ticker.ChangePrice), 2)),
			"(" + fmt.Sprintf("%.2f", pmChangeRate(ticker)*100) + "%)",
		}, " "),
		fmt.Sprintf(
			"[%s / 52주 신고가: %s(%s) / 52주 신저가: %s(%s)]",
			ticker.Market,
			humanize.CommafWithDigits(ticker.Highest52WeekPrice, 2),
			ticker.Highest52WeekDate,
			humanize.CommafWithDigits(ticker.Lowest52WeekPrice, 2),
			ticker.Lowest52WeekDate,
		),
	)
}
