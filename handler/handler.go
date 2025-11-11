package handler

import (
	"fmt"
	"log"
	"math"

	"github.com/Anywhay/weatherbot/clients/openweather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
	}
}

func (handler *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := handler.bot.GetUpdatesChan(u)

	for update := range updates {
		handler.handlerUpdate(update)
	}
}

func (handler *Handler) handlerUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	// If we got a message
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	coordinates, err := handler.owClient.Coordinates(update.Message.Text)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "не смогли получить координаты")
		msg.ReplyToMessageID = update.Message.MessageID
		handler.bot.Send(msg)
		return
	}

	weather, err := handler.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "не смогли получить прогноз погоды в этой местности")
		msg.ReplyToMessageID = update.Message.MessageID
		handler.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("Температура в %s: %d градусов", update.Message.Text, int(math.Round(weather.Temp))),
	)
	msg.ReplyToMessageID = update.Message.MessageID

	handler.bot.Send(msg)
}
