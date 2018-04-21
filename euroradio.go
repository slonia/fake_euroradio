package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"

	trans "github.com/aerokite/go-google-translate/pkg"
	"gopkg.in/telegram-bot-api.v4"
)

type Configuration struct {
	Token  string
	ChatId int64
}

var config Configuration

func main() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Configuration{}
	decoder.Decode(&config)
	argsWithoutProg := os.Args[1:]
	phrase := strings.Join(argsWithoutProg, " ")
	funLevel := analyze(phrase)
	belPhrase := translate(phrase)
	postResult(belPhrase, funLevel)
}

func analyze(phrase string) int {
	x := rand.Intn(5) - 5
	return x
}

func translate(phrase string) string {
	req := &trans.TranslateRequest{
		SourceLang: "ru",
		TargetLang: "be",
		Text:       phrase,
	}
	translated, err := trans.Translate(req)
	if err != nil {
		log.Fatalln(err)
	}
	return translated
}

func postResult(phrase string, level int) {
	result := phrase
	if level < 0 {
		result += strings.Repeat("(", -level)
	} else {
		result += strings.Repeat(")", level)
	}
	bot, _ := tgbotapi.NewBotAPI(config.Token)
	msg := tgbotapi.NewMessage(config.ChatId, result)
	bot.Send(msg)
}
