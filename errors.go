package httpclient

import (
	"fmt"
	"strings"
)

type ErrorTag string
type ErrorTagCollection []ErrorTag

func (c ErrorTagCollection) String(delimiter string) string {
	r := make([]string, len(c))
	for i, t := range c {
		r[i] = string(t)
	}
	return strings.Join(r, delimiter)
}

type BaseError struct {
	originalErr error
	tags        ErrorTagCollection
}

func (e *BaseError) Error() string {
	if len(e.tags) == 0 {
		return fmt.Sprintf("[httpclient] %s", e.originalErr.Error())
	}
	return fmt.Sprintf("[httpclient][%s] %s", e.tags.String("/"), e.originalErr.Error())
}

func (e *BaseError) Unwrap() error {
	return e.originalErr
}
