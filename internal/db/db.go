package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/fboucher/be-my-eyes/internal/models"
)

// DB represents the database connection
type DB struct {
	conn *sql.DB
}

// dbPath returns the path to the SQLite database file
func dbPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".config", "be-my-eyes")
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	
	return filepath.Join(configDir, "history.db"), nil
}

// Open opens the SQLite database and initializes the schema
func Open() (*DB, error) {
	path, err := dbPath()
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}

	// Initialize schema
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

// initSchema creates the necessary tables if they don't exist
func (db *DB) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS query_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		video_id TEXT NOT NULL,
		video_title TEXT NOT NULL,
		question TEXT NOT NULL,
		answer TEXT NOT NULL,
		error TEXT,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_query_history_video_id ON query_history(video_id);
	CREATE INDEX IF NOT EXISTS idx_query_history_created_at ON query_history(created_at);
	`

	_, err := db.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// SaveQuery saves a query and its result to the database
func (db *DB) SaveQuery(videoID, videoTitle, question, answer string, errMsg *string, status string) error {
	query := `
	INSERT INTO query_history (video_id, video_title, question, answer, error, status, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, videoID, videoTitle, question, answer, errMsg, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save query: %w", err)
	}

	return nil
}

// GetAllHistory retrieves all query history ordered by creation time (newest first)
func (db *DB) GetAllHistory() ([]models.QueryHistory, error) {
	query := `
	SELECT id, video_id, video_title, question, answer, error, status, created_at
	FROM query_history
	ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	var history []models.QueryHistory
	for rows.Next() {
		var h models.QueryHistory
		var errMsg sql.NullString

		if err := rows.Scan(&h.ID, &h.VideoID, &h.VideoTitle, &h.Question, &h.Answer, &errMsg, &h.Status, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if errMsg.Valid {
			h.Error = &errMsg.String
		}

		history = append(history, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return history, nil
}

// GetHistoryByVideoID retrieves query history for a specific video
func (db *DB) GetHistoryByVideoID(videoID string) ([]models.QueryHistory, error) {
	query := `
	SELECT id, video_id, video_title, question, answer, error, status, created_at
	FROM query_history
	WHERE video_id = ?
	ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	var history []models.QueryHistory
	for rows.Next() {
		var h models.QueryHistory
		var errMsg sql.NullString

		if err := rows.Scan(&h.ID, &h.VideoID, &h.VideoTitle, &h.Question, &h.Answer, &errMsg, &h.Status, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if errMsg.Valid {
			h.Error = &errMsg.String
		}

		history = append(history, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return history, nil
}
