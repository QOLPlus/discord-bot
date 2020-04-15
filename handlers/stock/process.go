package stock

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const Command = "주식"

func Process(s *discordgo.Session, m *discordgo.MessageCreate, data []string) {
	fmt.Println("잘못들어옴", data)
}
