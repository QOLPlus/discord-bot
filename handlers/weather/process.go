package weather

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"

	"github.com/QOLPlus/discord-bot/handlers"
)

var Registry = &handlers.HandlerRegistry{
	Commands: []string{"날씨", "웨더"},
	Proc:     Process,
}

func Process(s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	_, err := s.ChannelMessageSend(m.ChannelID, "날씨 기능은 아직 지원하지 않습니다. PR 환영\n요청한 지역:" + strings.Join(data, ","))
	if err != nil {
		fmt.Println("message 전송 실패!")
	}
}

