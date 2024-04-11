package openapi

type Response interface {
	Unmarshal([]byte) error
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int64  `json:"code"`
}
