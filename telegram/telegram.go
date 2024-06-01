package telegram

import (
	"context"
	"log"
	"math/rand"
	"taro-bot/gigachat"
	"taro-bot/tarot"

	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(ctx context.Context, update tgbotapi.Update, bot *tgbotapi.BotAPI, gigachatToken gigachat.Token, redisClient *redis.Client) {

	if update.Message != nil {

		// Чтение данных из Redis.
		result, err := redisClient.Get(ctx, "user_state-"+update.Message.From.UserName).Result()

		if result == "tarot_command-1" {
			msgText := Make3CardSpread(update, bot, gigachatToken, update.Message.Text)

			// Отправляем полученное сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}

			// Сбрасываем состояние пользователя
			err := redisClient.Set(ctx, "user_state-"+update.Message.From.UserName, "start", 0).Err()
			if err != nil {
				log.Fatalf("Ошибка записи в Redis: %v", err)
			}
		} else {
			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("3 карты", "tarot_command"),
				),
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите опцию:")
			msg.ReplyMarkup = inlineKeyboard
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}

	}

	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Request(callback); err != nil {
			log.Println(err)
		}
		if update.CallbackQuery.Data == "tarot_command" {

			err := redisClient.Set(ctx, "user_state-"+update.CallbackQuery.From.UserName, "tarot_command-1", 0).Err()
			if err != nil {
				log.Fatalf("Ошибка записи в Redis: %v", err)
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Напишите ваш вопрос:")
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func Make3CardSpread(update tgbotapi.Update, bot *tgbotapi.BotAPI, gigachatToken gigachat.Token, question string) string {
	// Создаем колоду Таро
	deck := tarot.CreateTarotDeck()
	var randomCards []tarot.TarotCard
	for i := 0; i < 3; i++ {
		randomIndex := rand.Intn(len(deck))
		randomCards = append(randomCards, deck[randomIndex])
	}

	var cards string
	for _, card := range randomCards {
		cards += card.Name + " (" + card.Suit + ")"
	}

	currentTime := time.Now()
	compareTime := currentTime.Add(-30 * time.Minute)

	if gigachatToken.Created.Before(compareTime) {
		var err error
		gigachatToken, err = gigachat.GetToken()
		if err != nil {
			log.Panic(err)
		}
	}
	message, err := gigachat.Generate(cards, question, gigachatToken.Value)
	if err != nil {
		log.Println(err)
	}
	return message

}
