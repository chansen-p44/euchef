package euchef

import (
	"log"
	"strings"

	"github.com/nlopes/slack"
)

// DeleteSlackMessage deletes the specified message
func DeleteSlackMessage(timestamp, channel, key string) (err error) {

	api := slack.New(key)

	_, _, err = api.DeleteMessage(channel, timestamp)
	if err != nil {
		log.Printf("api.DeleteMessage error: %s\n", err)
		return err
	}

	return nil
}

// PostSlackMessage post a message to slack
func PostSlackMessage(menu []MenuItem, channel, key string) (timestamp string, err error) {

	api := slack.New(key)

	attachments := []slack.Attachment{}

	for i, value := range menu {

		a := slack.Attachment{
			Title:    strings.Join(value.Title, "\n"),
			Text:     " ",
			ImageURL: value.ImageURL,
			ThumbURL: value.ImageURL,
		}

		// We want the first attachment to include the date of menu as pretext
		if i == 0 {
			a.Pretext = "Menu for " + value.Date
		}

		attachments = append(attachments, a)
	}

	msg := slack.PostMessageParameters{Attachments: attachments}

	var channelID string
	channelID, timestamp, err = api.PostMessage(channel, "", msg)
	if err != nil {
		log.Printf("api.PostMessage error: %s\n", err)
		return "", err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	return timestamp, err
}
