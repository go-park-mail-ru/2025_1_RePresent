package slot

type Format struct {
	Code        int    `json:"code" validate:"required"`
	Height      int    `json:"height" validate:"required,min=1"`
	Width       int    `json:"width" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1,max=255"`
}
