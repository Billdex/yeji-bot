package openapi

type Request interface {
	Method() string
	URI() string
	Marshal() ([]byte, error)
}
