package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
)

func startAutoReply(config *Configuration, rc *realtime.Client) {
	// wait group
	var wg sync.WaitGroup

	// only does text replies for now
	// add script based replies later
	wg.Add(len(config.Channels))
	for _, channelName := range config.Channels {
		go func(channelName string, rc *realtime.Client) {
			defer wg.Done()
			msgs := make(chan models.Message)
			channelID := channelIDs[channelName]
			err := rc.SubscribeToMessageStream(&models.Channel{ID: channelID}, msgs)

			if err != nil {
				log.Println(err)
				return
			}

			log.Println("Started listening on channel : ", channelName)

			for msg := range msgs {
				for regEx, replyAction := range RegexActions {
					if regEx.Match([]byte(msg.Msg)) {
						fmt.Println("Recieved msg details: ", msg) // TODO : remove later
						tempMsg := models.Message{}
						tempMsg.RoomID = msg.RoomID
						tempMsg.Msg = replyAction.textReply

						attachments := []models.Attachment{}

						if len(replyAction.imageURL) > 0 {
							attachments = append(attachments, models.Attachment{ImageURL: replyAction.imageURL})
						}
						if len(replyAction.videoURL) > 0 {
							attachments = append(attachments, models.Attachment{VideoURL: replyAction.videoURL})
						}
						if len(attachments) > 0 {
							tempMsg.Attachments = attachments
						}
						_, err := rc.SendMessage(&tempMsg)
						if err != nil {
							log.Println(err)
						}
						continue
					}
				}
			}
		}(channelName, rc)
	}

	wg.Wait()
}
