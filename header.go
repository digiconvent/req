package req

func (r *Request) Header(key, value string) *Request {
	r.headers[key] = value
	return r
}
