package req

import "net/http"

func Get(url string) *Request {
	return &Request{method: http.MethodGet, url: url, headers: map[string]string{}}
}
