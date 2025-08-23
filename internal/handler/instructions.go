package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"remnawave-tg-shop-bot/internal/config"
)

// InstructionsCallbackHandler shows the platform selection menu for available platforms
func (h Handler) InstructionsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery.Message.Message
	langCode := update.CallbackQuery.From.LanguageCode

	// Build dynamic keyboard based on configured URLs
	var buttons [][]models.InlineKeyboardButton
	var row []models.InlineKeyboardButton

	// Add buttons for platforms with configured URLs
	if config.InstructionsTVURL() != "" {
		row = append(row, models.InlineKeyboardButton{
			Text: h.translation.GetText(langCode, "instructions_tv_button"),
			URL:  config.InstructionsTVURL(),
		})
	}
	if config.InstructionsAndroidURL() != "" {
		row = append(row, models.InlineKeyboardButton{
			Text: h.translation.GetText(langCode, "instructions_android_button"),
			URL:  config.InstructionsAndroidURL(),
		})
	}

	// Add first row if it has buttons
	if len(row) > 0 {
		buttons = append(buttons, row)
	}

	// Create second row
	row = []models.InlineKeyboardButton{}
	if config.InstructionsIOSURL() != "" {
		row = append(row, models.InlineKeyboardButton{
			Text: h.translation.GetText(langCode, "instructions_ios_button"),
			URL:  config.InstructionsIOSURL(),
		})
	}
	if config.InstructionsPCURL() != "" {
		row = append(row, models.InlineKeyboardButton{
			Text: h.translation.GetText(langCode, "instructions_pc_button"),
			URL:  config.InstructionsPCURL(),
		})
	}

	// Add second row if it has buttons
	if len(row) > 0 {
		buttons = append(buttons, row)
	}

	// Add back button
	buttons = append(buttons, []models.InlineKeyboardButton{
		{Text: h.translation.GetText(langCode, "back_button"), CallbackData: CallbackStart},
	})

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    callback.Chat.ID,
		MessageID: callback.ID,
		Text:      h.translation.GetText(langCode, "instructions_menu_text"),
		ParseMode: models.ParseModeHTML,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: buttons},
	})
	if err != nil {
		slog.Error("Error sending instructions menu message", err)
	}
}

// The following handlers are no longer needed since we redirect to URLs,
// but keeping them for backward compatibility in case someone still uses callback-based approach

// InstructionsTVCallbackHandler redirects to TV instructions URL
func (h Handler) InstructionsTVCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This handler is kept for backward compatibility but shouldn't be called
	// since we now use direct URL buttons
	langCode := update.CallbackQuery.From.LanguageCode
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            h.translation.GetText(langCode, "instructions_tv_text"),
		ShowAlert:       true,
	})
	if err != nil {
		slog.Error("Error answering TV instructions callback", err)
	}
}

// InstructionsAndroidCallbackHandler redirects to Android instructions URL
func (h Handler) InstructionsAndroidCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This handler is kept for backward compatibility but shouldn't be called
	// since we now use direct URL buttons
	langCode := update.CallbackQuery.From.LanguageCode
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            h.translation.GetText(langCode, "instructions_android_text"),
		ShowAlert:       true,
	})
	if err != nil {
		slog.Error("Error answering Android instructions callback", err)
	}
}

// InstructionsIOSCallbackHandler redirects to iOS instructions URL
func (h Handler) InstructionsIOSCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This handler is kept for backward compatibility but shouldn't be called
	// since we now use direct URL buttons
	langCode := update.CallbackQuery.From.LanguageCode
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            h.translation.GetText(langCode, "instructions_ios_text"),
		ShowAlert:       true,
	})
	if err != nil {
		slog.Error("Error answering iOS instructions callback", err)
	}
}

// InstructionsPCCallbackHandler redirects to PC instructions URL
func (h Handler) InstructionsPCCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This handler is kept for backward compatibility but shouldn't be called
	// since we now use direct URL buttons
	langCode := update.CallbackQuery.From.LanguageCode
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            h.translation.GetText(langCode, "instructions_pc_text"),
		ShowAlert:       true,
	})
	if err != nil {
		slog.Error("Error answering PC instructions callback", err)
	}
}