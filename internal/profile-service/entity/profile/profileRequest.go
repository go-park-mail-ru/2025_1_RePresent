package profile

type ProfileRequest struct {
	Username    string `json:"username" validate:"required,min=5,max=50"`
	Description string `json:"description" validate:"max=200"`
}
