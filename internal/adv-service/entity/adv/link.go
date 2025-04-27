package adv

type Link struct {
	TextLink string
	UserID   int
	Height   int
	Width    int
}

type Slot struct {
	Link       string `json:"link" validate:"required"`
	UserID     int    `json:"user_id" validate:"required,gte=0"`
	FormatCode int    `json:"format_code" validate:"required,gte=0"`
}
