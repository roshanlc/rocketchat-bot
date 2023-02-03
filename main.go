package main

import "fmt"

func main() {
	config, err := readConfiguration()
	fmt.Println(config, err)
}

/*
func main() {
	server := url.URL{
		Scheme: "http",
		Host:   "localhost:3000",
	}

	realtimeClient, err := realtime.NewClient(&server, false)
	if err != nil {
		log.Fatalln(err)
	}

	loginDetails := models.UserCredentials{
		Email:    "bot@idk.com", //"haribdr@jankaritech.com",
		Password: "pass",
	}

	userDetails, err := realtimeClient.Login(&loginDetails)
	if err != nil {
		log.Fatalln(err)
	}
	defer realtimeClient.Close()

	restClient := rest.NewClient(&server, false)

	err = restClient.Login(&loginDetails)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("User details:\n", userDetails)

	channelName := "try"

	channelID, err := realtimeClient.GetChannelId(channelName)

	if err != nil {
		log.Println(err)
		return
	}

	channelDetails := models.Channel{
		ID: channelID,
	}
	fmt.Println("Channel Details:", channelDetails)

	err = realtimeClient.JoinChannel(channelID)

	if err != nil {
		log.Println(err)
		return
	}

	// channel for listening to realtime msgs
	msgs := make(chan models.Message)

	err = realtimeClient.SubscribeToMessageStream(&channelDetails, msgs)

	if err != nil {
		log.Println(err)
		return
	}

	reg1 := regexp.MustCompile(`(?mi)^.*\s*k.*c(h|hh|hhh)(a|aa|aaa)\s*$`)
	whoamiRegExp := regexp.MustCompile(`(?mi)^.*\s*who\s*are\s*you\s*$`)

	for msg := range msgs {
		if reg1.MatchString(msg.Msg) {
			tempMsg := models.Message{
				RoomID: msg.RoomID,
				Msg:    "K chaina sodhaana baru... list banara dinhu!!!!",
			}

			_, err := realtimeClient.SendMessage(&tempMsg)
			if err != nil {
				fmt.Println("Error while sending a message!")
				continue
			}
		} else if whoamiRegExp.MatchString(msg.Msg) {
			m := models.PostMessage{

				//RoomID: msg.RoomID,
				RoomID:  msg.RoomID,
				Channel: msg.Channel,
				// ParseUrls
				// Alias
				// Emoji
				// Avatar
				Attachments: []models.Attachment{{ImageURL: "https://nepstuff.com/wp-content/uploads/2020/07/ezgif.com-crop-2.gif"}},
			}
			restClient.PostMessage(&m)

			// restClient.Send(&channelDetails)
			// tempMsg := models.Message{
			// 	RoomID: msg.RoomID,
			// 	Msg:    "K chaina sodhaana baru... list banara dinhu!!!!",
			// }

			// _, err := realtimeClient.SendMessage(&tempMsg)
			// if err != nil {
			// 	fmt.Println("Error while sending a message!")
			// 	continue
			// }
		}

	}

}
*/
