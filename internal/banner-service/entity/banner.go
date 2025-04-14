package entity

type Banner struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"owner"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Balance     int    `json:"balance"`
	Link        string `json:"link"`
	Deleted     bool   `json:"deleted"`
}

var DefaultBanner = Banner{
	ID:          -1,
	OwnerID:     -1,
	Title:       "Здесь ничего нет",
	Content:     "8c7b890d0d478f3f880b5007781f8d2a",
	Description: "Но тут могла быть ваша реклама",
	Status:      0,
	Balance:     0,
	Link:        "http://re-target.ru",
	Deleted:     false,
}
