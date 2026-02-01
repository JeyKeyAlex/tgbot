package main

import (
	"net/http"
	"os"

	"log"

	"github.com/JeyKeyAlex/tgbot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(appConfig.Telegram.ApiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	if appConfig.Server.WebhookBaseURL != "" {
		runWebhook(bot, appConfig)
	} else {
		runLongPolling(bot)
	}
}

func runWebhook(bot *tgbotapi.BotAPI, appConfig *config.Configuration) {
	port := appConfig.Server.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	webhookURL := appConfig.Server.WebhookBaseURL + "/webhook"
	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = bot.Request(wh); err != nil {
		log.Fatal(err)
	}
	log.Printf("Webhook set: %s", webhookURL)

	updates := bot.ListenForWebhook("/webhook")

	// Health check для Render (проверка живости сервиса)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	})

	go func() {
		log.Printf("Listening on :%s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	for update := range updates {
		handleUpdate(bot, update)
	}
}

func runLongPolling(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handleUpdate(bot, update)
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я очень люблю свою Катечку")
	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Send error: %v", err)
	}
}
