package req

func (e *RequestError) Unwrap() error {
	return e.Err
}
