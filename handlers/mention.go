package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`https?://discord(app)?.com/channels/(?P<serverId>\w+)/(?P<channelId>\w+)/(?P<messageId>\w+)`)

type Mention struct {
	Content   string
	ServerId  string
	ChannelId string
	MessageId string
}

func parseMatched(datum []string) map[string]string {
	parsed := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i == 0 {
			parsed["content"] = datum[i]
			continue
		}

		if len(name) == 0 {
			continue
		}

		parsed[name] = datum[i]
	}
	return parsed
}

func fetchMentions(content string) []Mention {
	var result []Mention

	match := re.FindAllStringSubmatch(content, -1)
	if len(match) == 0 {
		return result
	}

	for _, datum := range match {
		parsed := parseMatched(datum)
		result = append(result, Mention{
			Content: parsed["content"],
			ServerId: parsed["serverId"],
			ChannelId: parsed["channelId"],
			MessageId: parsed["messageId"],
		})
	}

	return result
}

func mentionHandled(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	mentions := fetchMentions(m.Content)
	if len(mentions) == 0 {
		return false
	}

	handled, err := handleMentions(s, m, mentions)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return handled
}

func handleMentions(s *discordgo.Session, m *discordgo.MessageCreate, mentions []Mention) (bool, error) {
	var embeds []*discordgo.MessageEmbed

	content := m.Content

	for i, mention := range mentions {
		originalMessage, err := s.ChannelMessage(mention.ChannelId, mention.MessageId)
		if err != nil {
			continue
		}

		channel, err := s.Channel(mention.ChannelId)
		if err != nil {
			continue
		}

		var messageImage *discordgo.MessageEmbedImage
		if len(originalMessage.Attachments) > 0 {
			attachment := originalMessage.Attachments[0]
			messageImage = &discordgo.MessageEmbedImage{
				URL: attachment.URL,
				ProxyURL: attachment.ProxyURL,
				Width: attachment.Width,
				Height: attachment.Height,
			}
		}

		description := fmt.Sprintf("%s [🔗](%s)", originalMessage.Content, mention.Content)

		embeds = append(embeds, &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name: originalMessage.Author.Username,
				IconURL: originalMessage.Author.AvatarURL(""),
			},
			Description: description,
			Timestamp: string(originalMessage.Timestamp),
			Image: messageImage,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "#" + channel.Name,
			},
		})

		content = strings.ReplaceAll(content, mention.Content,fmt.Sprintf("️️||📎%d||", i+1))
	}

	webhook, err := getOrCreateWebhook(s, m)
	if err != nil {
		return false, err
	}

	_, err = s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Content:   content,
		Username:  m.Author.Username,
		AvatarURL: m.Author.AvatarURL(""),
		Embeds:    embeds,
	})
	if err != nil {
		return false, err
	}

	// message delete 는 실패해도 크리티컬 하지 않다.
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	return true, nil
}


func getOrCreateWebhook(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.Webhook, error) {
	webhooks, err := s.ChannelWebhooks(m.ChannelID)
	if err != nil {
		return nil, err
	}

	var qolplusWebhook *discordgo.Webhook
	for _, webhook := range webhooks {
		if webhook.Name == fmt.Sprintf("QOLPlus_%s", m.ChannelID) {
			qolplusWebhook = webhook
			break
		}
	}

	if qolplusWebhook == nil {
		qolplusWebhook, err = s.WebhookCreate(m.ChannelID, fmt.Sprintf("QOLPlus_%s", m.ChannelID), "")
		if err != nil {
			return nil, err
		}
	}

	return qolplusWebhook, nil
}
