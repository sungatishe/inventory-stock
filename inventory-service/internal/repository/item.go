package repository

import (
	"database/sql"
	"errors"
	"inventory-service/internal/models"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type ItemRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) AddItem(item *models.Item) error {
	query := "INSERT INTO items (name, quantity, price, user_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, item.Name, item.Quantity, item.Price, item.UserID).Scan(&item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepository) UpdateItem(item *models.Item) error {
	query := "UPDATE items SET name = $1, quantity = $2, price = $3 WHERE id = $4 AND user_id = $5"
	result, err := r.db.Exec(query, item.Name, item.Quantity, item.Price, item.ID, item.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrItemNotFound
	}
	return nil
}

func (r *ItemRepository) GetItemByID(itemID, userID string) (*models.Item, error) {
	query := "SELECT id, name, quantity, price, user_id from items WHERE id = $1 AND user_id = $2"
	row := r.db.QueryRow(query, itemID, userID)

	item := &models.Item{}
	err := row.Scan(&item.ID, &item.Name, &item.Quantity, &item.Price, &item.UserID)

	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) DeleteItem(itemID, userID string) error {
	query := "DELETE FROM items WHERE id = $1 AND user_id = $2"
	result, err := r.db.Exec(query, itemID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrItemNotFound
	}

	return nil
}

func (r *ItemRepository) GetAllItems(userID string) ([]*models.Item, error) {
	query := "SELECT id, name, quantity, price, user_id FROM items WHERE user_id = $1"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Price, &item.UserID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
