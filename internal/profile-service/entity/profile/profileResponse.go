package profile

type ProfileResponse struct {
	Username    string `json:"username" validate:"required,min=5,max=35"`
	Email       string `json:"email" validate:"required,email"`
	Description string `json:"description" validate:"min=0,max=200"`
	Balance     int    `json:"balance" validate:"gte=0"`
	Role        int    `json:"role" validate:"required,gte=1,lte=2"`
}
