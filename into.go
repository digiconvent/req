package req

import "encoding/json"

// v must be a pointer to a struct with json tags.
func (r *Request) Into(v any) error {
	body, _, err := r.Do()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
