package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// BroadcastCommandHandler handles the /broadcast command for sending messages to all users
func (h Handler) BroadcastCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Extract the message to broadcast from the command
	messageText := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/broadcast"))
	
	if messageText == "" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Please provide a message to broadcast.\n\nUsage: /broadcast <message>",
		})
		if err != nil {
			slog.Error("Error sending broadcast usage message", err)
		}
		return
	}

	// Get all customers from database
	customers, err := h.customerRepository.FindAll(ctx)
	if err != nil {
		slog.Error("Error retrieving customers for broadcast", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "❌ Error retrieving users from database",
		})
		if err != nil {
			slog.Error("Error sending error message", err)
		}
		return
	}

	if len(customers) == 0 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "ℹ️ No users found in database",
		})
		if err != nil {
			slog.Error("Error sending no users message", err)
		}
		return
	}

	// Send initial status message
	statusMsg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("📤 Starting broadcast to %d users...", len(customers)),
	})
	if err != nil {
		slog.Error("Error sending initial status message", err)
		return
	}

	// Broadcast to all users with progress tracking
	successCount := 0
	failCount := 0
	totalUsers := len(customers)

	for i, customer := range customers {
		// Send message to customer
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    customer.TelegramID,
			Text:      messageText,
			ParseMode: models.ParseModeHTML,
		})

		if err != nil {
			slog.Error("Failed to send broadcast message", 
				"telegramId", customer.TelegramID, 
				"error", err)
			failCount++
		} else {
			successCount++
		}

		// Update progress every 10 users or at the end
		if (i+1)%10 == 0 || i == totalUsers-1 {
			progressText := fmt.Sprintf(
				"📤 Broadcasting progress:\n"+
					"• Total users: %d\n"+
					"• Sent: %d\n"+
					"• Failed: %d\n"+
					"• Progress: %d/%d (%.1f%%)",
				totalUsers,
				successCount,
				failCount,
				i+1,
				totalUsers,
				float64(i+1)/float64(totalUsers)*100,
			)

			_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    statusMsg.Chat.ID,
				MessageID: statusMsg.ID,
				Text:      progressText,
			})
			if err != nil {
				slog.Error("Error updating progress message", err)
			}
		}

		// Add small delay to avoid hitting rate limits
		time.Sleep(50 * time.Millisecond)
	}

	// Send final result
	finalText := fmt.Sprintf(
		"✅ Broadcast completed!\n\n"+
			"📊 Results:\n"+
			"• Total users: %d\n"+
			"• Successfully sent: %d\n"+
			"• Failed: %d\n"+
			"• Success rate: %.1f%%",
		totalUsers,
		successCount,
		failCount,
		float64(successCount)/float64(totalUsers)*100,
	)

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    statusMsg.Chat.ID,
		MessageID: statusMsg.ID,
		Text:      finalText,
	})
	if err != nil {
		slog.Error("Error sending final broadcast result", err)
	}

	slog.Info("Broadcast completed", 
		"total", totalUsers,
		"success", successCount, 
		"failed", failCount,
		"admin", update.Message.From.ID)
}