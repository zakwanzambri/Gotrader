package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/zakwanzambri/Gotrader/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection
func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables creates the necessary database tables
func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS trading_signals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		action TEXT NOT NULL,
		price REAL NOT NULL,
		entry_price REAL,
		stop_loss REAL,
		take_profit REAL,
		confidence TEXT NOT NULL,
		source TEXT NOT NULL,
		raw_message TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		status TEXT NOT NULL DEFAULT 'ACTIVE'
	);

	CREATE INDEX IF NOT EXISTS idx_symbol ON trading_signals(symbol);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON trading_signals(timestamp);
	CREATE INDEX IF NOT EXISTS idx_status ON trading_signals(status);
	`

	_, err := db.conn.Exec(query)
	return err
}

// InsertSignal inserts a new trading signal
func (db *DB) InsertSignal(signal *models.TradingSignal) error {
	query := `
	INSERT INTO trading_signals (symbol, action, price, entry_price, stop_loss, take_profit, confidence, source, raw_message, timestamp, status)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.conn.Exec(query,
		signal.Symbol, signal.Action, signal.Price, signal.EntryPrice,
		signal.StopLoss, signal.TakeProfit, signal.Confidence, signal.Source,
		signal.RawMessage, signal.Timestamp, signal.Status)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	signal.ID = int(id)
	return nil
}

// GetSignals retrieves signals with pagination and filtering
func (db *DB) GetSignals(limit, offset int, symbol, status string) ([]models.TradingSignal, error) {
	query := `SELECT id, symbol, action, price, entry_price, stop_loss, take_profit, confidence, source, raw_message, timestamp, status
			  FROM trading_signals WHERE 1=1`
	args := []interface{}{}

	if symbol != "" {
		query += " AND symbol = ?"
		args = append(args, symbol)
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []models.TradingSignal
	for rows.Next() {
		var signal models.TradingSignal
		err := rows.Scan(&signal.ID, &signal.Symbol, &signal.Action, &signal.Price,
			&signal.EntryPrice, &signal.StopLoss, &signal.TakeProfit, &signal.Confidence,
			&signal.Source, &signal.RawMessage, &signal.Timestamp, &signal.Status)
		if err != nil {
			log.Printf("Error scanning signal: %v", err)
			continue
		}
		signals = append(signals, signal)
	}

	return signals, nil
}

// GetSignalCount returns the total count of signals with optional filtering
func (db *DB) GetSignalCount(symbol, status string) (int, error) {
	query := "SELECT COUNT(*) FROM trading_signals WHERE 1=1"
	args := []interface{}{}

	if symbol != "" {
		query += " AND symbol = ?"
		args = append(args, symbol)
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	var count int
	err := db.conn.QueryRow(query, args...).Scan(&count)
	return count, err
}

// UpdateSignalStatus updates the status of a signal
func (db *DB) UpdateSignalStatus(id int, status string) error {
	query := "UPDATE trading_signals SET status = ? WHERE id = ?"
	_, err := db.conn.Exec(query, status, id)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}