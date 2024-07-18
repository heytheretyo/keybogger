package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"
)

type BufferStats struct {
	Keystrokes   int
	LeftClicks   int
	RightClicks  int
	MiddleClicks int
	MouseTravel  float64
	ScrollWheels int
}

func initialize_database() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home directory: %v", err)
	}
	keyloggerDir := filepath.Join(homeDir, ".keylogger")

	err = os.MkdirAll(keyloggerDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create .keylogger directory: %v", err)
	}

	dbPath := filepath.Join(keyloggerDir, "events.db")

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS event_counts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		left_clicks INTEGER DEFAULT 0,
		right_clicks INTEGER DEFAULT 0,
		middle_clicks INTEGER DEFAULT 0,
		keystrokes INTEGER DEFAULT 0,
		mouse_travel_distance REAL DEFAULT 0,
		scroll_wheel_movements INTEGER DEFAULT 0,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO event_counts (left_clicks, right_clicks, middle_clicks, keystrokes, mouse_travel_distance, scroll_wheel_movements) VALUES (0, 0, 0, 0, 0, 0)`)
	if err != nil {
		log.Println("error initializing counts:", err)
	}

	bufferStats = load_daily_stats()
}

func save_bucket() {
	_, err := db.Exec(`INSERT INTO event_counts (left_clicks, right_clicks, middle_clicks, keystrokes, mouse_travel_distance, scroll_wheel_movements) VALUES (?, ?, ?, ?, ?, ?)`,
		bucket.LeftClicks,
		bucket.RightClicks,
		bucket.MiddleClicks,
		bucket.Keystrokes,
		bucket.MouseTravel,
		bucket.ScrollWheels,
	)
	if err != nil {
		log.Println("error saving bucket:", err)
	}
	bucket = EventBucket{} // reset the bucket
}

func load_daily_stats() BufferStats {
	// Get the start and end of the current day
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
	SELECT
		COALESCE(SUM(left_clicks), 0),
		COALESCE(SUM(right_clicks), 0),
		COALESCE(SUM(middle_clicks), 0),
		COALESCE(SUM(keystrokes), 0),
		COALESCE(SUM(mouse_travel_distance), 0),
		COALESCE(SUM(scroll_wheel_movements), 0)
	FROM event_counts
	WHERE timestamp BETWEEN ? AND ?`

	row := db.QueryRow(query, startOfDay, endOfDay)
	var stats BufferStats

	err := row.Scan(&stats.LeftClicks, &stats.RightClicks, &stats.MiddleClicks,
		&stats.Keystrokes, &stats.MouseTravel, &stats.ScrollWheels)

	if err != nil {
		log.Printf("error reading daily stats: %v", err)
	}

	return stats
}
