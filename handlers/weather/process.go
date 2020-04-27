package weather

import (
	"fmt"
	"github.com/QOLPlus/core/commands/weather"
	"github.com/QOLPlus/discord-bot/refs"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var Registry = &refs.HandlerRegistry{
	Commands: []string{"날씨", "웨더", "온도", "기온"},
	Proc:     Process,
}

func Process(s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	for _, keyword := range data {
		if len(keyword) == 0 {
			keyword = "반포동"
		}
		keyword = strings.ReplaceAll(keyword, " ", "")
		data, err := weather.FetchLocationByKeyword(keyword)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "알 수 없는 에러가 발생했습니다.")
			return
		}

		region := data.GetFirstRegionCode()
		if region == nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "찾을 수 없는 지역입니다.")
			return
		}
		fetched, err := weather.FetchWeather(region)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "알 수 없는 에러가 발생했습니다.")
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, buildMessage(fetched))
	}
}

func buildMessage(r *weather.FetchWeatherResult) string {
	beforeString := "어제와 같음"
	before := int(r.DiffWithYesterday)
	if before != 0 {
		beforeString = fmt.Sprintf("어제대비 %d", before)
	}

	return fmt.Sprintf(
		"**%s**: %s **%d**°C (체감 %d°, %s°) [최저 %d° ~ 최고 %d°]",
		r.Location,
		r.Status,

		int(r.Temperature),
		int(r.TemperatureDayFeel),
		beforeString,
		int(r.TemperatureDayLow),
		int(r.TemperatureDayHigh),
	)
}
