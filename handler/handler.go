package handler

import (
	"fmt"
	"log"
	"math"

	"github.com/Anywhay/weatherbot/clients/openweather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot        *tgbotapi.BotAPI
	owClient   *openweather.OpenWeatherClient
	userCities map[int64]string
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:        bot,
		owClient:   owClient,
		userCities: make(map[int64]string),
	}
}

func (handler *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := handler.bot.GetUpdatesChan(u)

	for update := range updates {
		handler.HandlerUpdate(update)
	}
}

func (handler *Handler) HandlerUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "city":
			city := update.Message.CommandArguments()
			handler.userCities[update.Message.Chat.ID] = city
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Город %s сохранен", city))
			handler.bot.Send(msg)
			msg.ReplyToMessageID = update.Message.MessageID
			return

		case "weather":
			city, ok := handler.userCities[update.Message.From.ID]
			if !ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Сначала установите город с помощью команды /city"))
				handler.bot.Send(msg)
				msg.ReplyToMessageID = update.Message.MessageID
				return
			}
			coordinates, err := handler.owClient.Coordinates(city)
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
				fmt.Sprintf("Температура в городе %s: %d градусов", city, int(math.Round(weather.Temp))),
			)
			msg.ReplyToMessageID = update.Message.MessageID

			handler.bot.Send(msg)
			return

		default:
			log.Printf("New comand [%s], %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "такой команды не существует")
			msg.ReplyToMessageID = update.Message.MessageID
			handler.bot.Send(msg)
			return
		}
	}
}
