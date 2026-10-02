package req

import (
	"io"
	"net/http"
)

func Post(url string, body io.Reader) *Request {
	return &Request{method: http.MethodPost, url: url, body: body, headers: map[string]string{}}
}
