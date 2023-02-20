package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
)

// Result holds output from scripts
type Result struct {
	Output string
	Err    error
}

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

						if !replyAction.useScript {
							// send reply msg
							tempMsg := models.Message{}
							tempMsg.RoomID = msg.RoomID
							tempMsg.Msg = replyAction.textReply

							attachments := []models.Attachment{}
							// add attachments if there are any
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
						} else {
							// buffered channel
							reply := make(chan Result, 1)

							// exec the script
							go execExternalScript(reply, replyAction.scriptLocation, true, msg.User.Name, msg.RoomID, msg.Msg)

							// get result of execution
							result := <-reply
							// close channel
							close(reply)

							if result.Err != nil {
								log.Println(result.Err)
								continue
							}

							// parse the json from script
							msgReply := MsgReply{}
							err = json.Unmarshal([]byte(result.Output), &msgReply)
							if err != nil {
								log.Println(err)
								continue
							}

							// send reply message to channel
							tempMsg := models.Message{}
							tempMsg.RoomID = msg.RoomID
							tempMsg.Msg = msgReply.TextReply

							attachments := []models.Attachment{}

							if len(msgReply.ImageURL) > 0 {
								attachments = append(attachments, models.Attachment{ImageURL: msgReply.ImageURL})
							}
							if len(msgReply.VideoURL) > 0 {
								attachments = append(attachments, models.Attachment{VideoURL: msgReply.VideoURL})
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
			}
		}(channelName, rc)
	}

	wg.Wait()
}

func execExternalScript(reply chan Result, scriptLocation string, inputData bool, userName, roomID, message string) {
	// ticker
	ticker := time.NewTicker(1 * time.Second)
	count := 0

	result := Result{}
	go func() {
		var data []byte
		var err error
		if inputData {
			json := fmt.Sprintf("{\"user_name\" : \"%v\",\"room_id\": \"%s\", \"msg\" : \"%s\"}", userName, roomID, message)
			// execute the script
			data, err = exec.Command("node", scriptLocation, "-d", json).CombinedOutput()
		} else {
			// execute the script
			data, err = exec.Command("node", scriptLocation).CombinedOutput()
		}

		// set the output from script onto the Result entity
		if err != nil {
			result = Result{
				Output: "",
				Err:    err,
			}
		} else {
			result = Result{
				Output: string(data),
				Err:    nil,
			}
		}
	}()

	// loop over ticker
	for range ticker.C {
		// increment seconds count
		count++
		// if seconds count has crossed limit
		if count > 9 {
			reply <- Result{
				Output: "",
				Err:    fmt.Errorf("time out for while waiting for result"),
			}
			ticker.Stop()
			break
		}
		// incase of result, send the result to the channel
		if result.Output != "" || result.Err != nil {
			reply <- result
			ticker.Stop()
			break
		}
	}
}
