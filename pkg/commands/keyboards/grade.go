package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// The Color keyboard shows the range of color that we support in the app
var Grade = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5A", "5A"),
		tgbotapi.NewInlineKeyboardButtonData("5B", "5B"),
		tgbotapi.NewInlineKeyboardButtonData("5C", "5C"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("6A", "6A"),
		tgbotapi.NewInlineKeyboardButtonData("6B", "6B"),
		tgbotapi.NewInlineKeyboardButtonData("6C", "6C"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("7A", "7A"),
		tgbotapi.NewInlineKeyboardButtonData("7B", "7B"),
		tgbotapi.NewInlineKeyboardButtonData("7C", "7C"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("8A", "8A"),
		tgbotapi.NewInlineKeyboardButtonData("8B", "8B"),
		tgbotapi.NewInlineKeyboardButtonData("8C", "8C"),
	),
)
