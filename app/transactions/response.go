package transactions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Error struct {
		Code      int    `json:"code"`
		Status    string `json:"status"`
		Message   string `json:"message"`
		RequestID string `json:"requestId"`
	} `json:"error"`
}

func NewResponse(status, message string, data interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(code int, message, status, requestId string) *ErrorResponse {
	errorResp := &ErrorResponse{}
	errorResp.Error.Code = code
	errorResp.Error.Status = status
	errorResp.Error.Message = message
	errorResp.Error.RequestID = requestId
	return errorResp
}

func writeJSONResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Check the type of the response and encode accordingly
	switch v := response.(type) {
	case *Response:
		// Success response
		if err := json.NewEncoder(w).Encode(v); err != nil {
			http.Error(w, "Failed to encode success response", http.StatusInternalServerError)
			log.Printf("Error encoding success response: %v", err)
		}
	case *ErrorResponse:
		// Error response
		if err := json.NewEncoder(w).Encode(v); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			log.Printf("Error encoding error response: %v", err)
		}
	default:
		// Default case for unsupported response types
		http.Error(w, "Unsupported response type", http.StatusInternalServerError)
		log.Printf("Unsupported response type: %T", v)
	}
}

func GenerateNewSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response := NewResponse(fmt.Sprintf("%d", http.StatusOK), "Success", data)
	writeJSONResponse(w, response, http.StatusOK)
}

func GenerateNewErrorResponse(w http.ResponseWriter, status, message, requestId string, code int) {
	response := NewErrorResponse(code, message, status, requestId)
	writeJSONResponse(w, response, http.StatusBadGateway)
}

func GenerateNewCreatedResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response := NewResponse(fmt.Sprintf("%d", http.StatusCreated), "Transaction created successfully", data)
	writeJSONResponse(w, response, http.StatusCreated)
}

func GenerateNewNoContentResponse() *Response {
	return NewResponse("204", "No content", nil)
}
