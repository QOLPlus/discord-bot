package handlers

import (
	"github.com/bwmarrin/discordgo"
	"strings"

	"github.com/QOLPlus/discord-bot/handlers/stock"
	"github.com/QOLPlus/discord-bot/handlers/weather"
)

type HandlerMap map[string]func(*discordgo.Session, *discordgo.MessageCreate, []string)
func (h HandlerMap) getKeys() []string {
	keys := make([]string, 0, len(h))
	for key := range h {
		keys = append(keys, key)
	}
	return keys
}

func HandlerFactory() interface{} {
	onHandlers := make(HandlerMap)
	onHandlers[stock.Command] = stock.Process
	onHandlers[weather.Command] = weather.Process
	onCommands := onHandlers.getKeys()

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		tasks := groupByCommands(parseTokens(m.Content), onCommands)
		if len(tasks) == 0 {
			return
		}

		for command, phrase := range tasks {
			onHandlers[command](s, m, phrase)
		}
	}
}

func parseTokens(content string) []string {
	var (
		sbOpen = false
		sbDepth = 0
		sbInnerContents = ""
	)
	var tokens []string

	for _, l := range content {
		letter := string(l)

		if !sbOpen {
			if letter != "[" {
				continue
			} else {
				sbDepth += 1
				sbOpen = true
			}
		} else {
			if letter == "]" {
				sbDepth -= 1
				if sbDepth == 0 {
					tokens = append(tokens, sbInnerContents)
					sbInnerContents = ""
					sbOpen = false
				} else {
					sbInnerContents += letter
				}
			} else if letter == "[" {
				sbDepth += 1
				sbInnerContents += letter
			} else {
				sbInnerContents += letter
			}
		}
	}

	return tokens
}

func groupByCommands(tokens []string, permits []string) map[string][]string {
	commandGroup := make(map[string][]string)
	permitMap := make(map[string]struct{}, len(permits))
	for _, s := range permits {
		permitMap[s] = struct{}{}
	}

	for _, token := range tokens {
		cmdWithPhrase := strings.SplitN(token, " ", 2)
		if len(cmdWithPhrase) < 2 {
			continue
		}

		command := cmdWithPhrase[0]
		phrase := cmdWithPhrase[1]

		if _, ok := permitMap[command]; ok {
			commandGroup[command] = append(commandGroup[command], phrase)
		}
	}
	return commandGroup
}
