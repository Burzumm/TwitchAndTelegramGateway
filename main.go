package main

import (
	"encoding/json"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	config := getConfiguration("./config.json")
	var telegramBotConnection = getTelegramBot(config)
	var twitchClient = twitch.NewClient(config.TwitchUserName, config.TwitchOuth)
	messageHandler(config, telegramBotConnection, twitchClient)
}

/*
Message handler.
*/
func messageHandler(config *Configuration, telegramBot *tgbotapi.BotAPI, twitchClient *twitch.Client) {
	go telegramMessageHandler(config, telegramBot, twitchClient)
	twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		var str string
		str += message.User.Name
		str += "\n"
		str += message.Message
		msg := tgbotapi.NewMessage(-1001662731743, str)
		telegramBot.Send(msg)
	})

	twitchClient.Join(config.TwitchChanel)
	err := twitchClient.Connect()
	if err != nil {
		panic(err)
	}
}

/*
Telegram message handler.
*/
func telegramMessageHandler(config *Configuration, telegramBot *tgbotapi.BotAPI, twitchClient *twitch.Client) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := telegramBot.GetUpdatesChan(u)
	for update := range updates {
		if update.ChannelPost != nil {
			var str string
			str += update.ChannelPost.Text
			twitchClient.Say(config.TwitchChanel, str)
		}
	}
}

/*
Returns Telegram Bot.
*/
func getTelegramBot(config *Configuration) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotApiKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

/*
Returns config app.
*/
func getConfiguration(path string) *Configuration {
	var m Configuration
	jsonFile, errorFile := os.Open(path)
	if errorFile != nil {
		log.Panic(errorFile)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	errJson := json.Unmarshal(byteValue, &m)
	if errJson != nil {
		log.Panic(errJson)
	}
	return &m
}

/*
Configuration App.
*/
type Configuration struct {
	TwitchUserName    string `json:"twitchUserName"`
	TwitchOuth        string `json:"twitchOuth"`
	TwitchChanel      string `json:"twitchChanel"`
	TelegramBotApiKey string `json:"telegramBotApiKey"`
	TelegramChatId    int64  `json:"telegramChatId"`
}
