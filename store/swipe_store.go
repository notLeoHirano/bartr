package store

import (
	"fmt"

	"github.com/notLeoHirano/bartr/models"
)

func (r *Store) CreateSwipe(swipe *models.Swipe) error {
	_, err := r.db.Exec(
		"INSERT INTO swipes (user_id, item_id, direction) VALUES (?, ?, ?)",
		swipe.UserID, swipe.ItemID, swipe.Direction,
	)
	return err
}

func (r *Store) UserSwipedRight(userID, itemID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM swipes 
		WHERE user_id = ? AND item_id = ? AND direction = 'right'
	`, userID, itemID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Matches

func (r *Store) CreateMatchIfNeeded(user1ID, user2ID, item1ID, item2ID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if user2 swiped right on item1
	var count int
	err = tx.QueryRow(`
		SELECT COUNT(*) 
		FROM swipes 
		WHERE user_id = ? AND item_id = ? AND direction = 'right'
	`, user2ID, item1ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}

	// Check if match already exists
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM matches 
		WHERE (user1_id = ? AND user2_id = ? AND item1_id = ? AND item2_id = ?)
		OR (user1_id = ? AND user2_id = ? AND item1_id = ? AND item2_id = ?)
	`, user1ID, user2ID, item1ID, item2ID, user2ID, user1ID, item2ID, item1ID).Scan(&count)

	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Create the match
	_, err = tx.Exec(
		"INSERT INTO matches (user1_id, user2_id, item1_id, item2_id) VALUES (?, ?, ?, ?)",
		user1ID, user2ID, item1ID, item2ID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Store) GetMatches(userID int) ([]models.MatchResponse, error) {
	query := `
		SELECT 
			m.id, m.user1_id, m.user2_id, m.item1_id, m.item2_id, m.created_at,
			i1.title, i2.title, u1.name, u2.name
		FROM matches m
		JOIN items i1 ON m.item1_id = i1.id
		JOIN items i2 ON m.item2_id = i2.id
		JOIN users u1 ON m.user1_id = u1.id
		JOIN users u2 ON m.user2_id = u2.id
		WHERE m.user1_id = ? OR m.user2_id = ?
		ORDER BY m.created_at DESC
	`

	rows, err := r.db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []models.MatchResponse{}
	for rows.Next() {
		var m models.MatchResponse
		if err := rows.Scan(&m.ID, &m.User1ID, &m.User2ID, &m.Item1ID, &m.Item2ID,
			&m.CreatedAt, &m.Item1Title, &m.Item2Title, &m.User1Name, &m.User2Name); err != nil {
			return nil, err
		}

		// Load comments for each match
		comments, err := r.GetComments(m.ID)
		if err == nil {
			m.Comments = comments
		}

		matches = append(matches, m)
	}

	return matches, rows.Err()
}

func (r *Store) MatchExists(user1ID, user2ID, item1ID, item2ID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM matches 
		WHERE (user1_id = ? AND user2_id = ? AND item1_id = ? AND item2_id = ?)
		OR (user1_id = ? AND user2_id = ? AND item1_id = ? AND item2_id = ?)
	`, user1ID, user2ID, item1ID, item2ID, user2ID, user1ID, item2ID, item1ID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking match existence: %w", err)
	}

	return count > 0, nil
}

func (r *Store) UserInMatch(matchID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM matches 
		WHERE id = ? AND (user1_id = ? OR user2_id = ?)
	`, matchID, userID, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}