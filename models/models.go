package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Item struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ItemWithOwner struct {
	Item
	OwnerName string `json:"owner_name"`
}

type Swipe struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ItemID    int       `json:"item_id" binding:"required"`
	Direction string    `json:"direction" binding:"required,oneof=left right"`
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

type MatchResponse struct {
	ID         int       `json:"id"`
	User1ID    int       `json:"user1_id"`
	User2ID    int       `json:"user2_id"`
	Item1ID    int       `json:"item1_id"`
	Item2ID    int       `json:"item2_id"`
	Item1Title string    `json:"item1_title"`
	Item2Title string    `json:"item2_title"`
	User1Name  string    `json:"user1_name"`
	User2Name  string    `json:"user2_name"`
	CreatedAt  time.Time `json:"created_at"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID        int       `json:"id"`
	MatchID   int       `json:"match_id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentRequest struct {
	MatchID int    `json:"match_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}