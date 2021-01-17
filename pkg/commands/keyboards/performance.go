package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// The Performance keyboard shows the range of performance results that can be
// achieved by a user when they attempted a route.
var Performance = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Flashed", "flashed"),
		tgbotapi.NewInlineKeyboardButtonData("Succeeded", "succeeded"),
		tgbotapi.NewInlineKeyboardButtonData("Failed", "failed"),
	),
)
