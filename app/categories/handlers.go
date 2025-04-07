package categories

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

// RegisterHandler registers the HTTP routes for categories
func (h *Handler) RegisterHandler(router *mux.Router) {
	// Route for GET request to retrieve all categories
	router.HandleFunc("/categories", h.GetCategories).Methods("GET")

	// Route for POST request to create a new category
	router.HandleFunc("/categories", h.CreateCategory).Methods("POST")

	// Route for GET request to retrieve a category by ID
	router.HandleFunc("/categories/{id}", h.GetCategoryByID).Methods("GET")

	// Route for PUT request to update a category by ID
	router.HandleFunc("/categories/{id}", h.UpdateCategory).Methods("PUT")

	// Route for DELETE request to delete a category by ID
	router.HandleFunc("/categories/{id}", h.DeleteCategory).Methods("DELETE")
}

// GetTransactionTypes handles GET requests to retrieve all transaction types
func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
		log.Printf("Error retrieving categories: %v", err)
		return
	}

	// Send the response wrapped in a success response format
	response := NewSuccessResponse(categories)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// CreateCategory handles POST requests to create a new category
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&category); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Call the service to create a new category
	createdCategory, err := h.service.CreateCategory(category)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		log.Printf("Error creating category: %v", err)
		return
	}

	// Wrap the response in a standard "created" response structure
	response := NewCreatedResponse([]Category{createdCategory})

	// Send JSON response with the created category
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// GetCategoryByID handles GET requests to retrieve a category by ID
func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to get the category by ID
	category, err := h.service.GetCategoryByID(idInt)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		log.Printf("Error retrieving category by ID: %v", err)
		return
	}

	// Wrap the response in a standard success response structure
	response := NewSuccessResponse([]Category{*category})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateCategory handles PUT requests to update an existing category by ID
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	var category Category
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&category); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Set the ID in the category (so the service knows which record to update)
	category.ID = idInt

	// Call the service to update the category
	updatedCategory, err := h.service.UpdateCategory(idInt, category)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		log.Printf("Error updating category: %v", err)
		return
	}

	// Wrap the response in a success response format
	response := NewSuccessResponse([]Category{updatedCategory})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteCategory handles DELETE requests to delete a category by ID
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	// Call the service to delete the transaction type
	err = h.service.DeleteCategory(idInt)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		log.Printf("Error deleting category: %v", err)
		return
	}

	// Send response indicating successful deletion
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
}
