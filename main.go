package main

import (
	"log"
	//"gopkg.in/telegram-bot-api.v4"
	"github.com/src/gopkg.in/telegram-bot-api.v4"
	"net/http"
	"strings"
	"encoding/json"
	"os"
	"io/ioutil"
)

type Config struct {
	Token string
}

func main() {
	file, _ := os.Open("config.json")

	contents,_ := ioutil.ReadFile("config.json")
	println(string(contents))

	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	//bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}


	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	//updates := bot.ListenForWebhook("/")

	var last string;

	for update := range updates {

		println("123")

		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID;
		title := update.Message.Chat.Title;

		println(chatId)
		println(title)

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		user, err := bot.GetMe()
		if err != nil {
			println(err)
		}

		if update.Message.Document != nil && "@" + user.UserName == last {
			fileId := update.Message.Document.FileID;
			println(fileId)
			file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileId});
			if err != nil {
				text := "Download exception"
				last = ""
				println(text)
				println(err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)
				continue
			}
			//println(file.FilePath)
			link := file.Link(config.Token)
			println(link);
			body := strings.NewReader(link)
			resp, err := http.Post("http://example.com/upload", "application/json", body)

			if err != nil {
				text := "Proxy doesn't respond"
				last = ""
				println(text)
				println(err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)
				continue
			}

			println(resp.Body)

			defer resp.Body.Close()

			var target string

			json.NewDecoder(resp.Body).Decode(target)

			var text = "Get response " + target
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
			last = "";
		} else {
			println("nothing")
		}
		last = update.Message.Text;
	}
}

