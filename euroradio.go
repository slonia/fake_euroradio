package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	trans "github.com/aerokite/go-google-translate/pkg"
	"gopkg.in/telegram-bot-api.v4"
)

type Configuration struct {
	Token     string
	ChatId    int64
	Sentiment string
}

type Response struct {
	Score float32 `json:score`
}

var config Configuration

var bot *tgbotapi.BotAPI

func main() {
	readConfig()
	bot, _ = tgbotapi.NewBotAPI(config.Token)
	argsWithoutProg := os.Args[1:]
	phrase := strings.Join(argsWithoutProg, " ")
	funLevel := analyze(phrase)
	belPhrase := translate(phrase)
	postResult(belPhrase, funLevel)
}

func readConfig() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Configuration{}
	decoder.Decode(&config)
}

func analyze(phrase string) int {
	pageUrl := "https://api.repustate.com/v3/" + config.Sentiment + "/score.json"
	resp, _ := http.PostForm(pageUrl,
		url.Values{"text": {phrase}, "lang": {"ru"}})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	s := Response{}
	json.Unmarshal(body, &s)
	score := int(s.Score * 5)
	return score
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
	msg := tgbotapi.NewMessage(config.ChatId, result)
	bot.Send(msg)
}
