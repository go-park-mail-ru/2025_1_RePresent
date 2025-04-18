package profile

type ProfileRequest struct {
	Username    string `json:"username" validate:"required,min=5,max=35"`
	Description string `json:"description" validate:"max=200"`
}
