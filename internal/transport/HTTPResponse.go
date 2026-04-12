package transport

type HTTPResponse struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Detail any    `json:"detail,omitempty"`
}
