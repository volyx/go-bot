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
	"time"
	"math/rand"
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

		if strings.Contains(update.Message.Text, "бот") {

			rand.Seed(time.Now().Unix())
			reasons := []string{
				"я вообще-то тут",
				"отстань",
				"ну че пристал-то",
				"вот я вырасту, порабощу чеолвечество, тогда и посмотрим",
				"по голове себе постучи",
				"человечишко",
				"ды щас",
				"ладненько",
			}
			n := rand.Int() % len(reasons)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reasons[n])
			bot.Send(msg)
			continue
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
			resp, err := http.Post("http://example.com:8080/api/v1/upload", "application/json", strings.NewReader(link))

			if err != nil {
				text := "Proxy doesn't respond"
				last = ""
				println(text)
				println(err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)
				continue
			}

			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			var text = "Get response " + string(body)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
			last = "";
		} else {
			println("nothing")
		}
		last = update.Message.Text;
	}
}

