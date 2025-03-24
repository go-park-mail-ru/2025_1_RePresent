package entity

type Banner struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"owner"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Balance     int    `json:"balance"`
}
