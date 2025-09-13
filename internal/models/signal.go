package models

import (
	"time"
)

// TradingSignal represents a trading signal captured from Telegram
type TradingSignal struct {
	ID          int       `json:"id" db:"id"`
	Symbol      string    `json:"symbol" db:"symbol"`
	Action      string    `json:"action" db:"action"` // BUY, SELL, HOLD
	Price       float64   `json:"price" db:"price"`
	EntryPrice  *float64  `json:"entry_price,omitempty" db:"entry_price"`
	StopLoss    *float64  `json:"stop_loss,omitempty" db:"stop_loss"`
	TakeProfit  *float64  `json:"take_profit,omitempty" db:"take_profit"`
	Confidence  string    `json:"confidence" db:"confidence"` // HIGH, MEDIUM, LOW
	Source      string    `json:"source" db:"source"`         // Telegram chat name/ID
	RawMessage  string    `json:"raw_message" db:"raw_message"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Status      string    `json:"status" db:"status"` // ACTIVE, CLOSED, EXPIRED
}

// SignalResponse represents the API response for signals
type SignalResponse struct {
	Signals []TradingSignal `json:"signals"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}

// WSMessage represents WebSocket message structure
type WSMessage struct {
	Type    string      `json:"type"`    // signal_update, new_signal, etc.
	Payload interface{} `json:"payload"`
}