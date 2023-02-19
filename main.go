package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sync"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/rest"
)

// reply actions for corresponding regex
type replyAction struct {
	useScript      bool
	scriptLocation string
	textReply      string
	imageURL       string
	videoURL       string
}

// Map compiled regex with their replies
// Since we'll be doing concurrent reads only,
// no mutex lock is necessary as of yet
// maybe aded later
var RegexActions = make(map[*regexp.Regexp]replyAction)

// Slice of auto running scripts
var AutoRunScripts []AutoRun

// map of channel names with their IDs
var channelIDs = make(map[string]string)

// clients
var realTimeClient *realtime.Client
var restClient *rest.Client

func main() {
	// Graceful shutdown from panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Following error has occurred:\n%v\n", r)
			log.Println("The system will not shutdown.")
		}
	}()
	config, err := readConfiguration()

	if err != nil {
		fmt.Println(err)
		return
	}
	// Print configuration details
	fmt.Println(config)

	err = config.checkValidity()
	if err != nil {
		log.Println(err)
		return
	}

	// login to server
	realTimeClient, restClient, err = config.login()

	// get Channel IDS
	channelIDs, err = config.getChannelIDs(realTimeClient)
	if err != nil {
		log.Println(err)
		return
	}

	// Join the channels
	for _, channelID := range channelIDs {
		realTimeClient.JoinChannel(channelID)
	}

	var wgMain sync.WaitGroup

	// Start the cronjobs (replying and auto run scripts)
	wgMain.Add(1)
	go func() {
		defer wgMain.Done()
		startAutoReply(config, realTimeClient)
	}()

	wgMain.Wait()
	// Close real time connection
	realTimeClient.Close()
}

// login to the server
// and return clients for realtime and rest
func (config *Configuration) login() (*realtime.Client, *rest.Client, error) {
	server := url.URL{
		Scheme: config.ServerDetails.Scheme,
		Host:   config.ServerDetails.URL,
	}
	loginDetails := models.UserCredentials{
		Email:    config.ServerDetails.Email,
		Password: config.ServerDetails.Password,
	}

	realtimeClient, err := realtime.NewClient(&server, false)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = realtimeClient.Login(&loginDetails)

	if err != nil {
		return nil, nil, err
	}

	restClient := rest.NewClient(&server, false)

	err = restClient.Login(&loginDetails)
	if err != nil {
		return nil, nil, err
	}

	return realtimeClient, restClient, nil
}

// return channel ids for all channels mentioned in the configuration
func (config *Configuration) getChannelIDs(rc *realtime.Client) (map[string]string, error) {
	var allChannels []string
	// add to chhanels
	allChannels = append(allChannels, config.Channels...)

	for _, val := range config.AutoRun {
		allChannels = append(allChannels, val.Channel)
	}

	var chanIDS = make(map[string]string)
	for _, val := range allChannels {
		if _, exists := chanIDS[val]; exists {
			continue
		}

		id, err := rc.GetChannelId(val)
		if err != nil {
			return nil, err
		}
		chanIDS[val] = id
	}
	return chanIDS, nil
}
