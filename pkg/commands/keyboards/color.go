package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// The Color keyboard shows the range of color that we support in the app
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

// The Color keyboard shows the range of color that we support in the app
var ColorChoices = []Choice{
	{Action: "red", Label: "Red"},
	{Action: "green", Label: "Green"},
	{Action: "blue", Label: "Blue"},
	{Action: "yellow", Label: "Yellow"},
	{Action: "orange", Label: "Orange"},
	{Action: "gray", Label: "Gray"},
}
