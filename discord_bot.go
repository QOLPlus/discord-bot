package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"

	"github.com/QOLPlus/discord-bot/handlers"
)

const gameTitle = "Quality of Life"

func main() {
	authToken := os.Getenv("DISCORD_AUTHENTICATION_TOKEN")
	bot, err := discordgo.New("Bot " + authToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Add handler
	bot.AddHandler(handlers.HandlerFactory())
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		_ = s.UpdateStatus(0, gameTitle)
	})

	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is running. To exit bot, press Ctrl-C.")
	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<- closeChan

	_ = bot.Close()
}
