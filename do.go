package req

import (
	"io"
	"net/http"
)

// returns the body, the response, and the stage of any failure.
func (r *Request) Do() ([]byte, *http.Response, error) {
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, nil, &RequestError{Stage: StageBuildRequest, Method: r.method, URL: r.url, Err: err}
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return nil, nil, &RequestError{Stage: StageSendRequest, Method: r.method, URL: r.url, Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, &RequestError{Stage: StageReadBody, Method: r.method, URL: r.url, Err: err}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, resp, &RequestError{Stage: StageStatus, Method: r.method, URL: r.url, Err: &StatusError{StatusCode: resp.StatusCode, Body: body}}
	}
	return body, resp, nil
}
