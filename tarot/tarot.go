package tarot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Структура для представления карты Таро
type TarotCard struct {
	Name   string
	Suit   string
	ImgURL tgbotapi.FileURL
	Number int
}
