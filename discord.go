package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func Discord(tasks chan Task) error {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return err
	}
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if s.State.User.ID == m.Author.ID {
			return
		}

		res := make(chan string)
		tasks <- Task{
			Source:   "discord",
			Content:  m.Content,
			Response: res,
		}

		s.ChannelMessageSend(m.ChannelID, <-res)

	})
	return discord.Open()
}
