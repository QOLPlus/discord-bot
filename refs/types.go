package refs

import "github.com/bwmarrin/discordgo"

type HandlerProc func(*discordgo.Session, *discordgo.MessageCreate, []string)
type HandlerMap map[string]HandlerProc
type HandlerRegistry struct {
	Commands []string
	Proc     HandlerProc
}
func (h HandlerMap) GetKeys() []string {
	keys := make([]string, 0, len(h))
	for key := range h {
		keys = append(keys, key)
	}
	return keys
}
func (h HandlerMap) Register(registry *HandlerRegistry) {
	for _, command := range registry.Commands {
		h[command] = registry.Proc
	}
}
