package database

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}

	// Set busy timeout to wait for locks
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		category TEXT,
		image_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS swipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		item_id INTEGER NOT NULL,
		direction TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (item_id) REFERENCES items(id),
		UNIQUE(user_id, item_id)
	);

	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user1_id INTEGER NOT NULL,
		user2_id INTEGER NOT NULL,
		item1_id INTEGER NOT NULL,
		item2_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user1_id) REFERENCES users(id),
		FOREIGN KEY (user2_id) REFERENCES users(id),
		FOREIGN KEY (item1_id) REFERENCES items(id),
		FOREIGN KEY (item2_id) REFERENCES items(id)
	);

	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		match_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE INDEX IF NOT EXISTS idx_comments_match_id ON comments(match_id);
	`

		if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Seed users
	if err := db.seedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Seed items
	if err := db.seedItems(); err != nil {
		return fmt.Errorf("failed to seed items: %w", err)
	}

	return nil
}

func (db *DB) seedUsers() error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		// Hash default passwords
		aliceHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		bobHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		charlieHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		_, err := db.Exec(`
			INSERT INTO users (name, email, password_hash) VALUES 
			('Alice', 'alice@example.com', ?),
			('Bob', 'bob@example.com', ?),
			('Charlie', 'charlie@example.com', ?)
		`, string(aliceHash), string(bobHash), string(charlieHash))
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) seedItems() error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		items := []struct {
			UserID      int
			Title       string
			Description string
			Category    string
			ImageURL    string
		}{
			{1, "Vintage Lamp", "A cozy lamp for your living room", "Home", ""},
			{1, "Mountain Bike", "Used bike in good condition", "Sports", ""},
			{2, "Cookbook", "Healthy recipes for beginners", "Books", ""},
			{2, "Acoustic Guitar", "Six-string guitar, lightly used", "Music", ""},
			{3, "Board Games Bundle", "Set of 5 popular games", "Toys & Games", ""},
			{3, "Desk Chair", "Ergonomic office chair", "Furniture", ""},
		}

		for _, item := range items {
			_, err := db.Exec(`
				INSERT INTO items (user_id, title, description, category, image_url)
				VALUES (?, ?, ?, ?, ?)
			`, item.UserID, item.Title, item.Description, item.Category, item.ImageURL)
			if err != nil {
				return err
			}
		}

	}

	return nil
}