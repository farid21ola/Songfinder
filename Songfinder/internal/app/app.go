package app

import (
	"Songfinder/config"
	"Songfinder/internal/app/bot"
	"path/filepath"
)

var privatePath = filepath.Join("private", "oauth.json")

func Run() {
	myBot := bot.BasicBot{
		ChannelName: config.ChannelName,
		MsgRate:     config.MsgRate,
		Name:        config.BotName,
		Port:        config.Port,
		PrivatePath: privatePath,
		Server:      config.Server,
	}
	myBot.Start()
}
