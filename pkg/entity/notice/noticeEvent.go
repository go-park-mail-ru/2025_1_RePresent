package notice

type NoticeEvent struct {
	UserID int     `json:"user_id"`
	Type   int     `json:"type"` // ex. low_balance, etc.
	Amount float64 `json:"amount,omitempty"`
}
