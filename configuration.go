// this file contains configuration structs
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// configuration file name
const configFile = "config.yml"

// Overall configuration
type Configuration struct {
	ServerDetails ServerDetails `yaml:"server_details"`
	Channels      []string      `yaml:"channels"` // Listens to these channels for msg replies
	ReplyTo       []ReplyTo     `yaml:"reply_to"` // Replies for corresponding msg
	AutoRun       []AutoRun     `yaml:"auto_run"` // Run script periodically
}

// ServerDetails struct contains
// server details and login credentials
type ServerDetails struct {
	Scheme   string `yaml:"scheme"`
	URL      string `yaml:"url"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

// convert to go lang regex, use regex101.com for this
// Supports nodejs script only
// If exec_script is provided then other later properties will be ignored
// the script should finish executing within 10 seconds
// the output from script should be on stdout and should follow the json structure
//
//	{
//	   "text_reply": "",
//	   "image_url": "",
//	   "video_url" : "",
//	}
type ReplyTo struct {
	MsgRegex   string `yaml:"msg_regex"`
	ExecScript string `yaml:"exec_script"`
	TextReply  string `yaml:"text_reply"`
	ImageURL   string `yaml:"image_url"`
	VideoURL   string `yaml:"video_url"`
}

// AutoRun struct holds values
// for scripts that run periodically
type AutoRun struct {
	ExecScript string `yaml:"exec_script"`
	Channel    string `yaml:"channel"`
}

// MsgReply struct is necessary parsing msg
// response from scripts
type MsgReply struct {
	TextReply string `json:"text_reply,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	VideoURL  string `json:"video_url,omitempty"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("Details:\n"+
		"\tScheme: %s\n"+
		"\tURL: %s\n"+
		"\tChannels To Listen For Reply: %v\n",
		c.ServerDetails.Scheme, c.ServerDetails.URL, c.Channels)
}

// readConfiguration reads configuration from
// config.yml
func readConfiguration() (*Configuration, error) {
	config := Configuration{}
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// checks if msg response sent from
// scripts is empty or not
func (m *MsgReply) isResponseEmpty() bool {
	if len(m.ImageURL) == 0 && len(m.VideoURL) == 0 && len(m.TextReply) == 0 {
		return true
	}
	return false
}

// checks if all fields of response
// are empty or not
// func (m *ReplyTo) isResponseEmpty() bool {
// 	if len(m.ImageURL) == 0 && len(m.VideoURL) == 0 && len(m.TextReply) == 0 {
// 		return true
// 	}
// 	return false
// }
