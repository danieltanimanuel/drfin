package accounts

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandler(router *mux.Router) {
	router.HandleFunc("/accounts", h.GetAllAcounts).Methods("GET")
	// router.HandleFunc("/accounts", h.CreateAccountGroup).Methods("POST")
	// router.HandleFunc("/accounts/{id}", h.GetAccountGroupByID).Methods("GET")
	// router.HandleFunc("/accounts/{id}", h.UpdateAccountGroup).Methods("PUT")
	// router.HandleFunc("/accounts/{id}", h.DeleteAccountGroup).Methods("DELETE")
}

// GetTransactionAccountGroups handles GET requests to retrieve all transaction accountgroups
func (h *Handler) GetAllAcounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.GetAllAccounts()
	if err != nil {
		http.Error(w, "Failed to retrieve AccountGroups", http.StatusInternalServerError)
		log.Printf("Error retrieving AccountGroups: %v", err)
		return
	}

	// Send the response wrapped in a success response format
	response := NewSuccessResponse(accounts)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// func (h *Handler) CreateAccountGroup(w http.ResponseWriter, r *http.Request) {
// 	var accountgroup AccountGroup
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&accountgroup); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		log.Printf("Error decoding request body: %v", err)
// 		return
// 	}

// 	createdAccountGroup, err := h.service.CreateAccountGroup(accountgroup)
// 	if err != nil {
// 		http.Error(w, "Failed to create accountgroup", http.StatusInternalServerError)
// 		log.Printf("Error creating accountgroup: %v", err)
// 		return
// 	}

// 	response := NewCreatedResponse([]AccountGroup{createdAccountGroup})

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		log.Printf("Error encoding response: %v", err)
// 	}
// }

// // GetAccountGroupByID handles GET requests to retrieve a accountgroup by ID
// func (h *Handler) GetAccountGroupByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	// Convert the ID from string to int
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		http.Error(w, "Invalid accountgroup ID", http.StatusBadRequest)
// 		log.Printf("Error converting ID to int: %v", err)
// 		return
// 	}

// 	// Call the service to get the accountgroup by ID
// 	accountgroups, err := h.service.GetAccountGroupByID(idInt)
// 	if err != nil {
// 		http.Error(w, "AccountGroup not found", http.StatusNotFound)
// 		log.Printf("Error retrieving accountgroup by ID: %v", err)
// 		return
// 	}

// 	// Wrap the response in a standard success response structure
// 	response := NewSuccessResponse([]AccountGroup{*accountgroups})

// 	// Send JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		log.Printf("Error encoding response: %v", err)
// 	}
// }

// // UpdateAccountGroup handles PUT requests to update an existing accountgroup by ID
// func (h *Handler) UpdateAccountGroup(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	// Convert the ID from string to int
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		http.Error(w, "Invalid accountgroup ID", http.StatusBadRequest)
// 		log.Printf("Error converting ID to int: %v", err)
// 		return
// 	}

// 	var accountgroups AccountGroup
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&accountgroups); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		log.Printf("Error decoding request body: %v", err)
// 		return
// 	}

// 	// Set the ID in the accountgroup (so the service knows which record to update)
// 	accountgroups.ID = idInt

// 	// Call the service to update the accountgroup
// 	updatedAccountGroup, err := h.service.UpdateAccountGroup(idInt, accountgroups)
// 	if err != nil {
// 		http.Error(w, "Failed to update accountgroup", http.StatusInternalServerError)
// 		log.Printf("Error updating accountgroup: %v", err)
// 		return
// 	}

// 	// Wrap the response in a success response format
// 	response := NewSuccessResponse([]AccountGroup{updatedAccountGroup})

// 	// Send JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		log.Printf("Error encoding response: %v", err)
// 	}
// }

// // DeleteAccountGroup handles DELETE requests to delete a accountgroup by ID
// func (h *Handler) DeleteAccountGroup(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	// Convert the ID from string to int
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		http.Error(w, "Invalid accountgroup ID", http.StatusBadRequest)
// 		log.Printf("Error converting ID to int: %v", err)
// 		return
// 	}

// 	// Call the service to delete the transaction accountgroup
// 	err = h.service.DeleteAccountGroup(idInt)
// 	if err != nil {
// 		http.Error(w, "Failed to delete accountgroup", http.StatusInternalServerError)
// 		log.Printf("Error deleting accountgroup: %v", err)
// 		return
// 	}

// 	// Send response indicating successful deletion
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
// }
