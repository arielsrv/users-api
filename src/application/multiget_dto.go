package application

type MultiGetDto[T any] struct {
	Code int `json:"code,omitempty"`
	Body T   `json:"body,omitempty"`
}
