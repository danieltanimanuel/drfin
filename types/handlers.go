package types

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

// NewHandler creates and returns a new Handler for transaction types
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterHandler registers the HTTP routes for Types
func (h *Handler) RegisterHandler(router *mux.Router) {
	// Route for GET request to retrieve all Types
	router.HandleFunc("/types", h.GetTypes).Methods("GET")

	// Route for POST request to create a new type
	router.HandleFunc("/types", h.CreateType).Methods("POST")

	// Route for GET request to retrieve a type by ID
	router.HandleFunc("/types/{id}", h.GetTypeByID).Methods("GET")

	// Route for PUT request to update a type by ID
	router.HandleFunc("/types/{id}", h.UpdateType).Methods("PUT")

	// Route for DELETE request to delete a type by ID
	router.HandleFunc("/types/{id}", h.DeleteType).Methods("DELETE")
}

// GetTransactionTypes handles GET requests to retrieve all transaction types
func (h *Handler) GetTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.service.GetAllTypes()
	if err != nil {
		http.Error(w, "Failed to retrieve Types", http.StatusInternalServerError)
		log.Printf("Error retrieving Types: %v", err)
		return
	}

	// Send the response wrapped in a success response format
	response := NewSuccessResponse(types)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// CreateType handles POST requests to create a new type
func (h *Handler) CreateType(w http.ResponseWriter, r *http.Request) {
	var types Type
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&types); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Call the service to create a new type
	createdType, err := h.service.CreateType(types)
	if err != nil {
		http.Error(w, "Failed to create type", http.StatusInternalServerError)
		log.Printf("Error creating type: %v", err)
		return
	}

	// Wrap the response in a standard "created" response structure
	response := NewCreatedResponse([]Type{createdType})

	// Send JSON response with the created type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// GetTypeByID handles GET requests to retrieve a type by ID
func (h *Handler) GetTypeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to get the type by ID
	types, err := h.service.GetTypeByID(idInt)
	if err != nil {
		http.Error(w, "Type not found", http.StatusNotFound)
		log.Printf("Error retrieving type by ID: %v", err)
		return
	}

	// Wrap the response in a standard success response structure
	response := NewSuccessResponse([]Type{*types})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateType handles PUT requests to update an existing type by ID
func (h *Handler) UpdateType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	var types Type
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&types); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Set the ID in the type (so the service knows which record to update)
	types.ID = idInt

	// Call the service to update the type
	updatedType, err := h.service.UpdateType(idInt, types)
	if err != nil {
		http.Error(w, "Failed to update type", http.StatusInternalServerError)
		log.Printf("Error updating type: %v", err)
		return
	}

	// Wrap the response in a success response format
	response := NewSuccessResponse([]Type{updatedType})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteType handles DELETE requests to delete a type by ID
func (h *Handler) DeleteType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to delete the transaction type
	err = h.service.DeleteType(idInt)
	if err != nil {
		http.Error(w, "Failed to delete type", http.StatusInternalServerError)
		log.Printf("Error deleting type: %v", err)
		return
	}

	// Send response indicating successful deletion
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
}
