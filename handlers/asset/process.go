package asset

import (
	"fmt"
	coreAsset "github.com/QOLPlus/core/commands/asset"
	"github.com/QOLPlus/discord-bot/refs"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"math"
	"regexp"
	"strings"
)


var commands = []string{"주식", "스톡", "스탁", "ST", "$"}

var Registry = &refs.HandlerRegistry{
	Commands: commands,
	Proc: Process,
}

var notFoundMessage string = fmt.Sprintf(
	"주식 정보를 찾을 수 없습니다. `[(%s) 검색어]`",
	strings.Join(commands, "|"),
)

const MaxCount int = 3

func Process(store *refs.Store, s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	assetStore := store.Asset
	assetStore.Reload(3600)

	var codes []string
	for _, phrase := range data {
		if len(phrase) == 0 {
			continue
		}

		var selectedAssets []coreAsset.Asset
		for _, datum := range *assetStore.Data {
			for _, asset := range *datum.Assets {
				if similar(asset, phrase) {
					selectedAssets = append(selectedAssets, asset)
				}
			}
		}
		if len(selectedAssets) == 0 {
			continue
		}

		var tCodes []string
		for _, asset := range selectedAssets {
			if include(asset, phrase) {
				tCodes = append(tCodes, asset.Code)
				break
			}
		}
		if len(tCodes) == 0 {
			tCodes = append(tCodes, selectedAssets[0].Code)
		}

		codes = append(codes, tCodes...)
	}

	if len(codes) == 0 {
		_, _ = s.ChannelMessage(m.ChannelID, notFoundMessage)
		return
	}

	securities, err := coreAsset.FetchSecuritiesByCodes(codes)
	if err != nil {
		_, _ = s.ChannelMessage(m.ChannelID, notFoundMessage)
		return
	}

	for idx, security := range *securities {
		if idx == MaxCount {
			warning := fmt.Sprintf("한번에 %d종목 이상을 조회할 수 없습니다.", MaxCount)
			_, _ = s.ChannelMessageSend(m.ChannelID, warning)
			break
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, buildSimpleMessage(security))
	}

}

func similar(asset coreAsset.Asset, keyword string) bool {
	match, _ := regexp.MatchString(
		keyword,
		strings.Join([]string{
			asset.KoreanName,
			asset.EnglishName,
			asset.TickerSymbol,
			asset.ShortCode,
			asset.Code,
		}, "|"),
	)
	return match
}

func include(asset coreAsset.Asset, keyword string) bool {
	for _, token := range []string{
		asset.KoreanName,
		asset.EnglishName,
		asset.TickerSymbol,
		asset.ShortCode,
		asset.Code,
	} {
		if token == keyword {
			return true
		}
	}

	return false
}

func buildSimpleMessage(security coreAsset.Security) string {
	return fmt.Sprintf(
		"%s **%s**: %s [%s]",
		flagOfCountry(security.ExchangeCountry),
		asterNameIfKosdaq(security),
		strings.Join([]string{
			humanize.Commaf(security.TradePrice),
			fmt.Sprintf("%s%.2f", arrowOfSecurity(security), math.Abs(security.ChangePrice)),
			"(" + fmt.Sprintf("%.2f", security.ChangeRate * 100) + "%)",
		}, " "),
		fmt.Sprintf(
			"%s / %s주 / %s",
			security.Market,
			humanize.Comma(security.AccTradeVolume),
			humanize.Commaf(security.GlobalAccTradePrice),
		),
	)
}


//tickerSymbol: AMZN
//shortCode: US.AMZN
//koreanName: 아마존닷컴
//englishName: AMAZON COM INC
//code: US.AMZN