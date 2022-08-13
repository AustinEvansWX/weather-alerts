package discord

import (
	"github.com/bwmarrin/discordgo"
)

func Login(token string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentGuildMessages

	err = session.Open()
	if err != nil {
		return nil, err
	}

	return session, nil
}
