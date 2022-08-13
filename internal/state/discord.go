package state

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type DiscordState struct {
	Mu      sync.RWMutex
	Session *discordgo.Session
}

var (
	discordState = DiscordState{}
)

func SetSession(session *discordgo.Session) {
	discordState.Session = session
}

func GetSession() *discordgo.Session {
	return discordState.Session
}
