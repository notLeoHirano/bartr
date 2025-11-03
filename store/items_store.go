package store

import "github.com/notLeoHirano/bartr/models"
func (r *Store) GetItems(userID int, excludeOwn bool) ([]models.ItemWithOwner, error) {
	query := `
		SELECT i.id, i.user_id, i.title, i.description, i.category, COALESCE(i.image_url, ''), i.created_at, u.name
		FROM items i
		JOIN users u ON i.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}

	if excludeOwn && userID > 0 {
		query += " AND i.user_id != ?"
		args = append(args, userID)
	}

	if userID > 0 {
		query += ` AND i.id NOT IN (
			SELECT item_id FROM swipes WHERE user_id = ?
		)`
		args = append(args, userID)
	}

	query += " ORDER BY i.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.ItemWithOwner{}
	for rows.Next() {
		var item models.ItemWithOwner
		if err := rows.Scan(&item.ID, &item.UserID, &item.Title, &item.Description,
			&item.Category, &item.ImageURL, &item.CreatedAt, &item.OwnerName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *Store) CreateItem(item *models.Item) error {
	result, err := r.db.Exec(
		"INSERT INTO items (user_id, title, description, category, image_url) VALUES (?, ?, ?, ?, ?)",
		item.UserID, item.Title, item.Description, item.Category, item.ImageURL,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = int(id)
	return nil
}

func (r *Store) DeleteItem(id int, userID int) (bool, error) {
	result, err := r.db.Exec("DELETE FROM items WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *Store) GetItemOwnerID(itemID int) (int, error) {
	var ownerID int
	err := r.db.QueryRow("SELECT user_id FROM items WHERE id = ?", itemID).Scan(&ownerID)
	return ownerID, err
}
