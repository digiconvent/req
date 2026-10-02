package req

type RequestError struct {
	Stage  Stage
	Method string
	URL    string
	Err    error
}
