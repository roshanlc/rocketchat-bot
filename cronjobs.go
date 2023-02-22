// this file contains cronjob related functions
// and also some internal cronjobs themselves
package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
)

func startCronJobs(rc *realtime.Client) {
	for _, scripts := range AutoRunScripts {
		go schedule(func() {
			reply := make(chan Result, 1)

			log.Printf("executing script %q at %v\n", scripts.ExecScript, time.Now())
			// exec the script
			go execExternalScript(reply, scripts.ExecScript, false, "", "", "")

			// get result of execution
			result := <-reply
			// close channel
			close(reply)

			if result.Err != nil {
				log.Println(result.Err)
				return
			}

			// parse the json from script
			msgReply := MsgReply{}
			err := json.Unmarshal([]byte(result.Output), &msgReply)
			if err != nil {
				log.Println(err)
				return
			}

			// incase of empty response, return from current execution
			if msgReply.isResponseEmpty() {
				log.Printf("empty response from script %q at %v\n", scripts.ExecScript, time.Now())
				return
			}
			// send reply message to channel
			tempMsg := models.Message{}
			tempMsg.RoomID = channelIDs[scripts.Channel]
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

			_, err = rc.SendMessage(&tempMsg)
			if err != nil {
				log.Println(err)
			}
		}, time.Duration(time.Minute*time.Duration(scripts.InEvery)))
	}
}

// schedule a cronjob
func schedule(f func(), interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			f()
		}
	}()
	return ticker
}
