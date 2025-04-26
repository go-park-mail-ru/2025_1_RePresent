package entity

type ResponseWithBody struct {
	Service ServiceResponse `json:"service"`
	Body    interface{}     `json:"body,omitempty"`
}

func NewResponseWithBody(isError bool, message string, body interface{}) ResponseWithBody {
	var service ServiceResponse

	if isError {
		service.Error = &message
	} else {
		service.Success = &message
	}

	return ResponseWithBody{
		Service: service,
		Body:    body,
	}
}
