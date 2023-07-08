package config

import "time"

const (
	//VK:
	VkID        = ""
	AccessToken = ""

	//Twitch:
	ChannelName = ""
	MsgRate     = time.Duration(20/30) * time.Millisecond
	BotName     = ""
	Port        = "6667"
	Server      = "irc.chat.twitch.tv"
)

//PrivatePath: "./private/oauth.json", - twitch oauth credentials
