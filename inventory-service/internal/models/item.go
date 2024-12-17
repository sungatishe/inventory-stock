package models

type Item struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Quantity int64   `json:"quantity"`
	Price    float32 `json:"price"`
	UserID   string  `json:"user_id"`
}
