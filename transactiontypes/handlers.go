package transactiontypes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handler struct is responsible for handling HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates and returns a new Handler for transaction transactionTypess
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterHandler registers the HTTP routes for TransactionTypes
func (h *Handler) RegisterHandler(router *mux.Router) {
	// Route for GET request to retrieve all TransactionTypes
	router.HandleFunc("/transactiontypes", h.GetTransactionTypes).Methods("GET")

	// Route for POST request to create a new type
	router.HandleFunc("/transactiontypes", h.CreateTransactionType).Methods("POST")

	// Route for GET request to retrieve a type by ID
	router.HandleFunc("/transactiontypes/{id}", h.GetTransactionTypeByID).Methods("GET")

	// Route for PUT request to update a type by ID
	router.HandleFunc("/transactiontypes/{id}", h.UpdateTransactionType).Methods("PUT")

	// Route for DELETE request to delete a type by ID
	router.HandleFunc("/transactiontypes/{id}", h.DeleteTransactionType).Methods("DELETE")
}

// GetTransactionTransactionTypes handles GET requests to retrieve all transaction transactionTypess
func (h *Handler) GetTransactionTypes(w http.ResponseWriter, r *http.Request) {
	transactionTypes, err := h.service.GetAllTransactionTypes()
	if err != nil {
		http.Error(w, "Failed to retrieve TransactionTypes", http.StatusInternalServerError)
		log.Printf("Error retrieving TransactionTypes: %v", err)
		return
	}

	// Send the response wrapped in a success response format
	response := NewSuccessResponse(transactionTypes)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// CreateTransactionType handles POST requests to create a new transactionType
func (h *Handler) CreateTransactionType(w http.ResponseWriter, r *http.Request) {
	var transactionType TransactionType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transactionType); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Call the service to create a new transactionType
	createdTransactionType, err := h.service.CreateTransactionType(transactionType)
	if err != nil {
		http.Error(w, "Failed to create transactionType", http.StatusInternalServerError)
		log.Printf("Error creating transactionType: %v", err)
		return
	}

	// Wrap the response in a standard "created" response structure
	response := NewCreatedResponse([]TransactionType{createdTransactionType})

	// Send JSON response with the created transactionType
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// GetTransactionTypeByID handles GET requests to retrieve a transactionType by ID
func (h *Handler) GetTransactionTypeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid transactionType ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to get the transactionType by ID
	transactionType, err := h.service.GetTransactionTypeByID(idInt)
	if err != nil {
		http.Error(w, "TransactionType not found", http.StatusNotFound)
		log.Printf("Error retrieving transactionType by ID: %v", err)
		return
	}

	// Wrap the response in a standard success response structure
	response := NewSuccessResponse([]TransactionType{*transactionType})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateTransactionType handles PUT requests to update an existing transactionType by ID
func (h *Handler) UpdateTransactionType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid transactionType ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	var transactionType TransactionType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transactionType); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Set the ID in the transactionType (so the service knows which record to update)
	transactionType.ID = idInt

	// Call the service to update the transactionType
	updatedTransactionType, err := h.service.UpdateTransactionType(idInt, transactionType)
	if err != nil {
		http.Error(w, "Failed to update transactionType", http.StatusInternalServerError)
		log.Printf("Error updating transactionType: %v", err)
		return
	}

	// Wrap the response in a success response format
	response := NewSuccessResponse([]TransactionType{updatedTransactionType})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteTransactionType handles DELETE requests to delete a transactionType by ID
func (h *Handler) DeleteTransactionType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid transactionType ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to delete the transaction transactionType
	err = h.service.DeleteTransactionType(idInt)
	if err != nil {
		http.Error(w, "Failed to delete transactionType", http.StatusInternalServerError)
		log.Printf("Error deleting transactionType: %v", err)
		return
	}

	// Send response indicating successful deletion
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
}
