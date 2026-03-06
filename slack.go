package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func Slack(tasks chan Task) error {
	api := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
	)

	client := socketmode.New(api)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					continue
				}
				client.Ack(*evt.Request)

				switch innerEvent := eventsAPIEvent.InnerEvent.Data.(type) {
				case *slackevents.MessageEvent:
					// ignore bot messages to avoid loops
					if innerEvent.BotID != "" || innerEvent.SubType == "bot_message" {
						continue
					}

					allowedUsers := os.Getenv("ALLOWED_SLACK_USERS")
					allowedChannels := os.Getenv("ALLOWED_SLACK_CHANNELS")

					if allowedUsers != "" {
						found := false
						for _, u := range splitsies(allowedUsers) {
							if innerEvent.User == u {
								found = true
								break
							}
						}
						if !found {
							continue
						}
					}

					text := innerEvent.Text
					channel := innerEvent.Channel

					mentioned := false
					mentionRe := regexp.MustCompile(`<@([A-Z0-9]+)>`)
					if m := mentionRe.FindStringSubmatch(text); len(m) > 1 && m[1] == os.Getenv("SLACK_APP_ID") {
						mentioned = true
						text = mentionRe.ReplaceAllString(text, "@Pinched")
					}

					if !mentioned {
						if allowedChannels == "" {
							// no channels configured, ignore non‑mentions
							continue
						}
						found := false
						for _, u := range splitsies(allowedChannels) {
							if channel == u {
								found = true
								break
							}
						}
						if !found {
							continue
						}
					}

					res := make(chan string)

					tasks <- Task{
						Source:   "slack",
						Content:  text,
						Response: res,
					}

					_, _, err := api.PostMessage(channel, slack.MsgOptionText(<-res, false))
					if err != nil {
						log.Printf("slack: failed to post message: %v", err)
					}

				case *slackevents.AppMentionEvent:
					text := innerEvent.Text
					channel := innerEvent.Channel
					text = regexp.MustCompile(`<@([A-Z0-9]+)>`).ReplaceAllString(text, "@Pinched")

					res := make(chan string)

					tasks <- Task{
						Source:   "slack",
						Content:  text,
						Response: res,
					}

					_, _, err := api.PostMessage(channel, slack.MsgOptionText(<-res, false))
					if err != nil {
						log.Printf("slack: failed to post message: %v", err)
					}
				}

			case socketmode.EventTypeConnecting:
				fmt.Println("slack: connecting...")
			case socketmode.EventTypeConnected:
				fmt.Println("slack: connected")
			case socketmode.EventTypeConnectionError:
				fmt.Println("slack: connection error")
			}
		}
	}()

	return client.Run()
}
