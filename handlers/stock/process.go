package stock

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"strings"

	"github.com/QOLPlus/core/commands/stock"
	"github.com/QOLPlus/discord-bot/refs"
)

var Registry = &refs.HandlerRegistry{
	Commands: []string{"주식", "스탁"},
	Proc:     Process,
}

const MaxCount int = 3

func Process(s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	codes := make([]string, 0)

	for _, keyword := range data {
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
		if len(searchCodes) == 1 {
			codes = append(codes, searchCodes[0])
		} else {
			codes = append(codes, searchCodes[0:2]...)
		}
	}

	if len(codes) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "주식 정보를 찾을 수 없습니다.")
		return
	}

	fetchSecuritiesResult, err := stock.FetchSecuritiesByCodes(codes)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "주식 정보를 불러올 수 없습니다.")
		return
	}

	for idx, security := range fetchSecuritiesResult.RecentSecurities {
		if idx == MaxCount {
			warning := fmt.Sprintf("한번에 %d종목 이상을 조회할 수 없습니다.", MaxCount)
			_, _ = s.ChannelMessageSend(m.ChannelID, warning)
			break
		}
		_, _ = s.ChannelMessageSendEmbed(
			m.ChannelID,
			buildEmbedFromSecurity(security),
		)
	}
}

func buildEmbedFromSecurity(security stock.FetchedSecurity) *discordgo.MessageEmbed {
	var (
		embedArrow string
		embedColor int
		embedMarket string
	)

	switch security.Change {
	case "UPPER_LIMIT":
		embedArrow = "↑"
		embedColor = 16062488
	case "RISE":
		embedArrow = "▲"
		embedColor = 16062488
	case "FALL":
		embedArrow = "▼"
		embedColor = 1794513
	case "LOWER_LIMIT":
		embedArrow = "↓"
		embedColor = 1794513
	default:
		embedArrow = ""
		embedColor = 0
	}

	if security.Market == "KOSDAQ" {
		embedMarket = "*"
	} else {
		embedMarket = ""
	}

	return &discordgo.MessageEmbed{
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://fn-chart.dunamu.com/images/kr/stock/d/mini/" + security.ShortCode + ".png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: fmt.Sprintf(
					"%s%s",
					security.Name,
					embedMarket,
				),
				Value: strings.Join([]string{
					humanize.Commaf(security.TradePrice),
					fmt.Sprintf("%s%.2f", embedArrow, security.ChangePrice),
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
