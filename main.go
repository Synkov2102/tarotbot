package main

import (
	"context"
	"log"
	"os"

	"github.com/Synkov2102/tarotbot/gigachat"
	"github.com/Synkov2102/tarotbot/postgres"
	"github.com/Synkov2102/tarotbot/redisConnector"
	"github.com/Synkov2102/tarotbot/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var (
	ctx = context.Background()
)

func main() {
	// Инициализация бота с токеном
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	postgres.InitDB()
	defer postgres.CloseDB()

	updates := bot.GetUpdatesChan(u)

	gigachatToken, err := gigachat.GetToken()
	if err != nil {
		log.Panic(err)
	}
	redisClient := redisConnector.GetRedisClient()
	// Обрабатываем входящие обновления
	for update := range updates {
		go telegram.HandleUpdate(ctx, update, bot, gigachatToken, redisClient)

	}
}
