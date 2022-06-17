package http_error

import (
	"fmt"
	"net/http"
)

// HTTPError extends the native error interface with some additional info for tracing and debugging.
type HTTPError struct {
	Headers      string
	ResponseBody string
	StatusCode   int
	URL          string
}

// NewHTTPError returns a new HTTPerror pointer.
func NewHTTPError(responseBody string, statusCode int, url string, headers http.Header) *HTTPError {
	return &HTTPError{
		Headers:      headersToString(headers),
		ResponseBody: responseBody,
		StatusCode:   statusCode,
		URL:          url,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[url: %s][error: Error calling API][status: %d][cause: %s]", e.URL, e.StatusCode, e.ResponseBody)
}

func headersToString(headers http.Header) string {
	result := ""
	for key, value := range headers {
		result = fmt.Sprintf("%s => %s,", key, value)
	}
	return result
}
