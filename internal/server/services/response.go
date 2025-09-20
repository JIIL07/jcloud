package services

import (
	"encoding/json"
	"net/http"
)

// ResponseService handles HTTP responses
type ResponseService struct{}

// NewResponseService creates a new response service
func NewResponseService() *ResponseService {
	return &ResponseService{}
}

// WriteJSON writes a JSON response
func (rs *ResponseService) WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response
func (rs *ResponseService) WriteError(w http.ResponseWriter, statusCode int, message string) error {
	errorResponse := map[string]string{
		"error":   http.StatusText(statusCode),
		"message": message,
	}
	return rs.WriteJSON(w, statusCode, errorResponse)
}

// WriteSuccess writes a success response
func (rs *ResponseService) WriteSuccess(w http.ResponseWriter, message string) error {
	successResponse := map[string]string{
		"message": message,
		"status":  "success",
	}
	return rs.WriteJSON(w, http.StatusOK, successResponse)
}

// WriteFile writes a file response
func (rs *ResponseService) WriteFile(w http.ResponseWriter, filename string, data []byte) error {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", string(rune(len(data))))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}

// WriteHTML writes an HTML response
func (rs *ResponseService) WriteHTML(w http.ResponseWriter, statusCode int, html string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(html))
	return err
}
