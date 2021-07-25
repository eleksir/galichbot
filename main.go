package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/onrik/micha"
	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("galichbot.ini")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	token := cfg.Section("").Key("BotApiKey").String()
	victim, err := cfg.Section("").Key("VictimID").Int64()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fileName := "phrases.txt"
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	speeches := strings.Split(string(fileBytes), "--")
	bot, err := micha.NewBot(token)
	if err != nil {
		log.Println(err)
		return
	}

	go bot.Start()
	rand.Seed(time.Now().Unix())

	for update := range bot.Updates() {
		if update.Message != nil {
			if update.Message.From.ID == victim {
				n := rand.Int() % len(speeches)
				var options = micha.SendMessageOptions{ReplyToMessageID: update.Message.MessageID}
				_, _ = bot.SendMessage(update.Message.Chat.ID, speeches[n], &options)
			}
		}
	}
}
