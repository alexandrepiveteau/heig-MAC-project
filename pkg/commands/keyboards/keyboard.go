package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// A sentinel value to display all the choices on a single line.
const SingleLine = 0

// A Choice is a mapping between a displayed value, and the action that will be
// actually sent to the bot.
type Choice struct {
	Label  string
	Action string
}

// GetActions returns all the possible actions for a list of choices keyboard
// choices.
func GetActions(choices []Choice) []string {
	actions := make([]string, 0)
	for _, choice := range choices {
		actions = append(actions, choice.Action)
	}
	return actions
}

// NewInlineKeyboard creates a new inline keyboard with a certain mapping, as
// well as a rowSize. You can use the SingleLine rowSize if you want to display
// all the choices on a single line.
func NewInlineKeyboard(
	choices []Choice,
	rowSize int,
) tgbotapi.InlineKeyboardMarkup {

	// Calculate the max row size.
	size := rowSize
	if rowSize <= SingleLine {
		size = len(choices)
	}

	rowsBuilder := make([][]tgbotapi.InlineKeyboardButton, 0)
	rowBuilder := make([]tgbotapi.InlineKeyboardButton, 0)
	rowIndex := 0
	for _, choice := range choices {

		// Append the choice to the current row.
		rowBuilder = append(rowBuilder, tgbotapi.NewInlineKeyboardButtonData(
			choice.Label,
			choice.Action,
		))
		rowIndex++

		// Append the newly built row.
		if rowIndex == size {
			rowsBuilder = append(rowsBuilder, rowBuilder)
			rowBuilder = make([]tgbotapi.InlineKeyboardButton, 0)
			rowIndex = 0
		}
	}

	// Build the actual keyboard as a combination of the rows.
	return tgbotapi.NewInlineKeyboardMarkup(rowsBuilder...)
}
