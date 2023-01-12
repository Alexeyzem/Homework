package main

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"strings"
	pc "taskbot/processing_command"
)

var (
	BotToken = "5897918940:AAGFWPiwuSIXAE4jyPdeR1Puk1KoL1kkACg"

	WebhookURL = "https://5f77-195-19-61-36.eu.ngrok.io"
)

type worker func(str string, user pc.MyUser) ([]pc.Out, error)

var commands = map[string]worker{
	"/tasks":    pc.Tasks,
	"/new":      pc.New,
	"/assign":   pc.Assign,
	"/unassign": pc.UnAssign,
	"/resolve":  pc.Resolve,
	"/my":       pc.My,
	"/owner":    pc.Owner,
}

func startTaskBot(ctx context.Context) error {
	commands["/help"] = func(str string, user pc.MyUser) ([]pc.Out, error) {
		text := "You can use:"
		for key := range commands {
			text += "\n" + key
		}
		var supportive = pc.Out{Message: text, ChatId: user.Id}
		var out []pc.Out
		out = append(out, supportive)
		return out, nil
	}
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return fmt.Errorf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		return fmt.Errorf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		return fmt.Errorf("SetWebhook failed: %s", err)
	}
	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is working"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)
	for update := range updates {
		var user = pc.MyUser{Name: update.Message.From.UserName, Id: update.Message.From.ID}
		key := update.Message.Text
		if string(key[0]) != `/` {
			continue
		}
		keys := strings.Split(key, " ")
		keys = strings.Split(keys[0], "_")
		funcWork, ok := commands[keys[0]]
		if !ok {
			wrongCommand := update.Message.Text + ": unknown command\nRun /help for usage."
			bot.Send(tgbotapi.NewMessage(update.Message.From.ID, wrongCommand))
			continue
		}
		key = strings.Replace(key, keys[0], "", 1)
		var msg []pc.Out
		if len(key) != 0 {
			msg, err = funcWork(key[1:], user)
		} else {
			msg, err = funcWork("", user)
		}
		if err != nil && !errors.Is(err, errors.New("internal error")) {
			bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "incorrect use of the command: "+keys[0]))
		} else if errors.Is(err, errors.New("internal error")) {
			bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Internal error, sorry.Try again later"))
		}
		for _, value := range msg {
			bot.Send(tgbotapi.NewMessage(value.ChatId, value.Message))
		}
	}
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
