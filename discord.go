package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Discord(tasks chan Task) error {
	fmt.Println("discord: connecting...")
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		return err
	}
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// in env: ALLOWED_DISCORD_USERS and ALLOWED_DISCORD_CHANNELS
		allowedUsers := os.Getenv("ALLOWED_DISCORD_USERS")
		allowedChannels := os.Getenv("ALLOWED_DISCORD_CHANNELS")

		if allowedUsers != "" {
			found := false
			for _, u := range splitsies(allowedUsers) {
				if m.Author.ID == u {
					found = true
					break
				}
			}
			if !found {
				return
			}
		}
		if allowedChannels != "" {
			found := false
			for _, u := range splitsies(allowedChannels) {
				if m.ChannelID == u {
					found = true
					break
				}
			}
			if !found {
				return
			}
		}

		res := make(chan string)

		err := s.ChannelTyping(m.ChannelID)
		if err != nil {
			println(err)
		}
		tasks <- Task{
			Source:   "discord",
			Content:  m.Content,
			Response: res,
		}

		s.ChannelMessageSend(m.ChannelID, <-res)

	})
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("discord: connected")
	})
	return discord.Open()
}
