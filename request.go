package req

import "io"

type Request struct {
	method  string
	url     string
	body    io.Reader
	headers map[string]string
}
