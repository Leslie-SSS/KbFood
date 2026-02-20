package dto

import (
	"encoding/json"
	"fmt"
	"net/http"

	"kbfood/internal/pkg/errors"
)

// Response represents a standard API response
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	Total int `json:"total,omitempty"`
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// Success creates a success response
func Success(data interface{}) Response {
	return Response{
		Code: http.StatusOK,
		Data: data,
	}
}

// SuccessWithMeta creates a success response with metadata
func SuccessWithMeta(data interface{}, meta Meta) Response {
	return Response{
		Code: http.StatusOK,
		Data: data,
		Meta: &meta,
	}
}

// Error creates an error response
func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// FromAppError creates an error response from AppError
func FromAppError(err *errors.AppError) Response {
	return Response{
		Code:    err.Code.HTTPStatus(),
		Message: err.Message,
	}
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, Error(code, message))
}

// WriteErr writes an error response from an error
func WriteErr(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		WriteJSON(w, FromAppError(appErr))
	} else {
		WriteError(w, http.StatusInternalServerError, "Internal server error")
	}
}
