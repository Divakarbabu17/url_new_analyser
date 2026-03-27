package utils

import (
    "fmt"
    "net/http"
)

// StatusText returns friendly message for a status code
func StatusText(code int) string {
    text := http.StatusText(code)
    if text == "" {
        text = "Unknown status"
    }
    return fmt.Sprintf("%d %s", code, text)
}

// IsStatusOK returns true if status code is 2xx
func IsStatusOK(code int) bool {
    return code >= 200 && code < 300
}

// FormatErrorMessage formats an error message with URL and HTTP status
func FormatErrorMessage(url string, code int, err error) string {
    if err != nil {
        return fmt.Sprintf("Error fetching %s: %v", url, err)
    }
    return fmt.Sprintf("Error fetching %s: %s", url, StatusText(code))
}