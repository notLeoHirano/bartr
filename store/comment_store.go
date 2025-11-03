package store

import "github.com/notLeoHirano/bartr/models"

func (r *Store) CreateComment(comment *models.Comment) error {
	result, err := r.db.Exec(
		"INSERT INTO comments (match_id, user_id, content) VALUES (?, ?, ?)",
		comment.MatchID, comment.UserID, comment.Content,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	comment.ID = int(id)
	return nil
}

func (r *Store) GetComments(matchID int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.match_id, c.user_id, u.name, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.match_id = ?
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.MatchID, &c.UserID, &c.UserName, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, rows.Err()
}