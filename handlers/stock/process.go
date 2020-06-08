package stock

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"math"
	"strings"
	"time"

	"github.com/QOLPlus/core/commands/stock"
	"github.com/QOLPlus/discord-bot/refs"
)

var Registry = &refs.HandlerRegistry{
	Commands: []string{"주식", "스탁"},
	Proc:     Process,
}

const MaxCount int = 3
const SubCmdDetail string = "상세"

type parsedConfigEntry struct {
	Keyword string
	Detail bool
}
type parsedConfig map[string]parsedConfigEntry

func Process(_ *refs.Store, s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	codes := make([]string, 0)
	config := parsedConfig{}

	for _, phrase := range data {
		if len(phrase) == 0 {
			continue
		}

		subCmdWithPhrase := strings.SplitN(phrase, " ", 2)
		var (
			subCmd  string
			keyword string
		)

		if len(subCmdWithPhrase) < 2 {
			subCmd = ""
			keyword = subCmdWithPhrase[0]

			// 그냥 지속
		} else {
			subCmd = subCmdWithPhrase[0]
			keyword = subCmdWithPhrase[1]
		}

		searchAssetsResult, err := stock.SearchAssetsByKeyword(keyword)
		if err != nil {
			continue
		}

		searchCodes := searchAssetsResult.FindExactlySameCodesByKeyword(keyword)
		if len(searchCodes) == 0 {
			searchCodes = searchAssetsResult.GetCodes()
			if len(searchCodes) == 0 {
				continue
			}
		}

		if len(searchCodes) > 3 {
			searchCodes = searchCodes[0:2]
		}

		for _, code := range searchCodes {
			config[code] = parsedConfigEntry{
				Keyword: keyword,
				Detail: subCmd == SubCmdDetail,
			}
			codes = append(codes, code)
		}
	}

	if len(codes) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "주식 정보를 찾을 수 없습니다. `[주식 검색어]` or `[주식 상세 검색어]`")
		return
	}

	fetchSecuritiesResult, err := stock.FetchSecuritiesByCodes(codes)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "주식 정보를 불러올 수 없습니다. `[주식 검색어]` or `[주식 상세 검색어]`")
		return
	}

	for idx, security := range fetchSecuritiesResult.RecentSecurities {
		if idx == MaxCount {
			warning := fmt.Sprintf("한번에 %d종목 이상을 조회할 수 없습니다.", MaxCount)
			_, _ = s.ChannelMessageSend(m.ChannelID, warning)
			break
		}

		doResponse(s, m, security, config)
	}
}

func doResponse(s *discordgo.Session, m *discordgo.MessageCreate, security stock.FetchedSecurity, config parsedConfig) {
	if config[security.Id].Detail {
		_, _ = s.ChannelMessageSendEmbed(
			m.ChannelID,
			buildDetailEmbedMessage(security),
		)
	} else {
		_, _ = s.ChannelMessageSend(
			m.ChannelID,
			buildSimpleMessage(security),
		)
	}
}

func buildSimpleMessage(security stock.FetchedSecurity) string {
	return fmt.Sprintf(
		"**%s**: %s [%s]",
		asterNameIfKosdaq(security),
		strings.Join([]string{
			humanize.Commaf(security.TradePrice),
			fmt.Sprintf("%s%.2f", arrowOfSecurity(security), math.Abs(security.ChangePrice)),
			"(" + fmt.Sprintf("%.2f", security.ChangePriceRate * 100) + "%)",
		}, " "),
		fmt.Sprintf(
			"%s주 / %s원",
			humanize.Comma(security.AccTradeVolume),
			humanize.Commaf(security.GlobalAccTradePrice),
		),
	)
}

func buildDetailEmbedMessage(security stock.FetchedSecurity) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: colorOfSecurity(security),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf(
				"https://fn-chart.dunamu.com/images/kr/stock/d/mini/%s.png?%d",
				security.ShortCode,
				time.Now().Unix(),
			),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: asterNameIfKosdaq(security),
				Value: strings.Join([]string{
					humanize.Commaf(security.TradePrice),
					fmt.Sprintf("%s%.2f", arrowOfSecurity(security), math.Abs(security.ChangePrice)),
					"(" + fmt.Sprintf("%.2f", security.ChangePriceRate * 100) + "%)",
				}, " "),
			},
			{
				Name: "시가 / 종가",
				Value: fmt.Sprintf("%s / %s", humanize.Commaf(security.OpeningPrice), humanize.Commaf(security.PrevClosingPrice)),
				Inline: true,
			},
			{
				Name: "고가 / 저가",
				Value: fmt.Sprintf("%s / %s", humanize.Commaf(security.HighPrice), humanize.Commaf(security.LowPrice)),
				Inline: true,
			},
			{
				Name: "거래량",
				Value: fmt.Sprintf(
					"%s 주 (%s 원)",
					humanize.Comma(security.AccTradeVolume),
					humanize.Commaf(security.GlobalAccTradePrice),
				),
			},
		},
	}
}
