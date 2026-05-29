package api

import "fmt"

// APIError represents a non-2xx response from the backend.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api %d: %s", e.Status, e.Message)
}

// IsNotFound reports whether the error is a 404 from the backend.
func IsNotFound(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.Status == 404
	}
	return false
}

// IsRateLimited reports whether the error is a 429.
func IsRateLimited(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.Status == 429
	}
	return false
}

// IsClientError reports whether the error is a 4xx (excluding 404/429).
func IsClientError(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.Status >= 400 && e.Status < 500
	}
	return false
}
