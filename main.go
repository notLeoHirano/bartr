package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Item struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

type Swipe struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ItemID    int       `json:"item_id"`
	Direction string    `json:"direction"` // "left" or "right"
	CreatedAt time.Time `json:"created_at"`
}

type Match struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	Item1ID   int       `json:"item1_id"`
	Item2ID   int       `json:"item2_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ItemWithOwner struct {
	Item
	OwnerName string `json:"owner_name"`
}

var db *sql.DB

func main() {

}