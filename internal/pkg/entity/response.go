package entity

type ServiceResponse struct {
	Error   *string `json:"error,omitempty"`   // Может быть null или отсутствовать
	Success *string `json:"success,omitempty"` // Может быть null или отсутствовать
}

type Response struct {
	Service ServiceResponse `json:"service"`
}

func NewResponse(isError bool, message string) Response {
	var service ServiceResponse

	if isError {
		service.Error = &message
	} else {
		service.Success = &message
	}

	return Response{
		Service: service,
	}
}
