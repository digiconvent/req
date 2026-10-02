package req

import "fmt"

func (e *StatusError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, string(e.Body))
}
