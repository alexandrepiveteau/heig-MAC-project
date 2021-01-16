package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var Color = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Red", "red"),
		tgbotapi.NewInlineKeyboardButtonData("Green", "green"),
		tgbotapi.NewInlineKeyboardButtonData("Blue", "blue"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Yellow", "yellow"),
		tgbotapi.NewInlineKeyboardButtonData("Orange", "orange"),
		tgbotapi.NewInlineKeyboardButtonData("Gray", "gray"),
	),
)
