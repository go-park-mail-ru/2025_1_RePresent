package profile

import "gopkg.in/inf.v0"

type ProfileResponse struct {
	ID          int     `json: "id"`
	Username    string  `json:"username" validate:"required,min=5,max=50"`
	Email       string  `json:"email" validate:"required,email"`
	Description string  `json:"description" validate:"min=0,max=200"`
	Balance     inf.Dec `json:"balance"`
	Role        int     `json:"role" validate:"required,gte=1,lte=2"`
}
