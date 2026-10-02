package req

type StatusError struct {
	StatusCode int
	Body       []byte
}
