package req

import "fmt"

func (e *RequestError) Error() string {
	return fmt.Sprintf("%s %s: %s: %v", e.Method, e.URL, e.Stage, e.Err)
}
