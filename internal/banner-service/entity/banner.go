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
	Content:     "0d93f25e17810c6fc8666718f37a5203", //хардкоднуто
	Description: "Но тут могла быть ваша реклама",
	Status:      0,
	Balance:     0,
	Link:        "http://re-target.ru",
	Deleted:     false,
}
