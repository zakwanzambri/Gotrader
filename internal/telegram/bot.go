package telegram

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zakwanzambri/Gotrader/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	chatIDs  []int64
	onSignal func(*models.TradingSignal)
}

// NewBot creates a new Telegram bot instance
func NewBot(token string, chatIDs []int64, onSignal func(*models.TradingSignal)) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		api:      bot,
		chatIDs:  chatIDs,
		onSignal: onSignal,
	}, nil
}

// Start starts the bot to listen for messages
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Check if message is from monitored chat
		if !b.isMonitoredChat(update.Message.Chat.ID) {
			continue
		}

		signal := b.parseMessage(update.Message)
		if signal != nil {
			signal.Source = update.Message.Chat.Title
			if signal.Source == "" {
				signal.Source = update.Message.Chat.UserName
			}
			if signal.Source == "" {
				signal.Source = strconv.FormatInt(update.Message.Chat.ID, 10)
			}

			if b.onSignal != nil {
				b.onSignal(signal)
			}
		}
	}
}

// isMonitoredChat checks if the chat ID is in the monitored list
func (b *Bot) isMonitoredChat(chatID int64) bool {
	for _, id := range b.chatIDs {
		if id == chatID {
			return true
		}
	}
	return false
}

// parseMessage attempts to parse trading signals from message text
func (b *Bot) parseMessage(message *tgbotapi.Message) *models.TradingSignal {
	text := strings.ToUpper(strings.TrimSpace(message.Text))
	
	if text == "" {
		return nil
	}

	// Common trading signal patterns
	patterns := []struct {
		regex  *regexp.Regexp
		action string
	}{
		{regexp.MustCompile(`(BUY|LONG)\s+([A-Z]{3,10})(?:/[A-Z]{3,4})?\s+(?:@|AT)?\s*([0-9.]+)`), "BUY"},
		{regexp.MustCompile(`(SELL|SHORT)\s+([A-Z]{3,10})(?:/[A-Z]{3,4})?\s+(?:@|AT)?\s*([0-9.]+)`), "SELL"},
		{regexp.MustCompile(`([A-Z]{3,10})(?:/[A-Z]{3,4})?\s+(BUY|LONG)\s+(?:@|AT)?\s*([0-9.]+)`), "BUY"},
		{regexp.MustCompile(`([A-Z]{3,10})(?:/[A-Z]{3,4})?\s+(SELL|SHORT)\s+(?:@|AT)?\s*([0-9.]+)`), "SELL"},
	}

	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(text)
		if len(matches) >= 3 {
			signal := &models.TradingSignal{
				RawMessage: message.Text,
				Timestamp:  time.Unix(int64(message.Date), 0),
				Status:     "ACTIVE",
				Action:     pattern.action,
			}

			// Extract symbol and price based on pattern
			if pattern.action == "BUY" || pattern.action == "SELL" {
				if strings.Contains(matches[1], "BUY") || strings.Contains(matches[1], "SELL") || strings.Contains(matches[1], "LONG") || strings.Contains(matches[1], "SHORT") {
					signal.Symbol = matches[2]
					if price, err := strconv.ParseFloat(matches[3], 64); err == nil {
						signal.Price = price
						signal.EntryPrice = &price
					}
				} else {
					signal.Symbol = matches[1]
					if price, err := strconv.ParseFloat(matches[3], 64); err == nil {
						signal.Price = price
						signal.EntryPrice = &price
					}
				}
			}

			// Extract additional information
			b.extractAdditionalInfo(signal, text)
			
			// Set confidence based on keywords
			signal.Confidence = b.determineConfidence(text)

			return signal
		}
	}

	return nil
}

// extractAdditionalInfo extracts stop loss, take profit, etc.
func (b *Bot) extractAdditionalInfo(signal *models.TradingSignal, text string) {
	// Stop Loss patterns
	slPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:SL|STOP\s*LOSS|STOP)[\s:@]*([0-9.]+)`),
		regexp.MustCompile(`STOP[\s:@]*([0-9.]+)`),
	}

	for _, pattern := range slPatterns {
		if matches := pattern.FindStringSubmatch(text); len(matches) > 1 {
			if sl, err := strconv.ParseFloat(matches[1], 64); err == nil {
				signal.StopLoss = &sl
				break
			}
		}
	}

	// Take Profit patterns
	tpPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:TP|TAKE\s*PROFIT|TARGET)[\s:@]*([0-9.]+)`),
		regexp.MustCompile(`TARGET[\s:@]*([0-9.]+)`),
	}

	for _, pattern := range tpPatterns {
		if matches := pattern.FindStringSubmatch(text); len(matches) > 1 {
			if tp, err := strconv.ParseFloat(matches[1], 64); err == nil {
				signal.TakeProfit = &tp
				break
			}
		}
	}
}

// determineConfidence determines signal confidence based on keywords
func (b *Bot) determineConfidence(text string) string {
	highConfidenceKeywords := []string{"STRONG", "HIGH", "CONFIDENT", "SURE", "CONFIRMED"}
	lowConfidenceKeywords := []string{"WEAK", "LOW", "UNCERTAIN", "MAYBE", "POSSIBLE"}

	for _, keyword := range highConfidenceKeywords {
		if strings.Contains(text, keyword) {
			return "HIGH"
		}
	}

	for _, keyword := range lowConfidenceKeywords {
		if strings.Contains(text, keyword) {
			return "LOW"
		}
	}

	return "MEDIUM"
}