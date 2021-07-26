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
		log.Fatalln(err)
	}

	token := cfg.Section("").Key("BotApiKey").String()
	victim, err := cfg.Section("").Key("VictimID").Int64()
	if err != nil {
		log.Fatalln(err)
	}

	fileName := "phrases.txt"
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	speeches := strings.Split(string(fileBytes), "--")
	bot, err := micha.NewBot(token)
	if err != nil {
		log.Fatalln(err)
	}

	go bot.Start()
	rand.Seed(time.Now().Unix())

	for update := range bot.Updates() {
		if update.Message != nil {
			log.Printf(
					"%s %s %s (%d) from %s (%d) says: $s",
					update.Message.From.Username,
					update.Message.From.FirstName,
					update.Message.From.LastName,
					update.Message.From.ID,
					update.Message.Chat.Title,
					update.Message.Chat.ID,
					update.Message.Text
			)

			if update.Message.From.ID == victim {
				n := rand.Int() % len(speeches)
				var options = micha.SendMessageOptions{ReplyToMessageID: update.Message.MessageID}
				_, err = bot.SendMessage(update.Message.Chat.ID, speeches[n], &options)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("Bot answers: %s", speeches[n])
				}
			}
		}
	}
}
