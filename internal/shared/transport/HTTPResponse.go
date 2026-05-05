package transport

type HTTPResponse struct {
	status  int
	Message any `json:"message"`
}

func NewHTTPResponse(status int, message any) HTTPResponse {
	return HTTPResponse{status: status, Message: message}
}

func (r HTTPResponse) Status() int {
	return r.status
}
