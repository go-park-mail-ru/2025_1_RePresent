package entity

type Banner struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"url_content"`
	Status      int    `json:"status"`
}
