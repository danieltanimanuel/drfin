package transactions

import (
	"encoding/json"
	"fmt"
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

func (h *Handler) RegisterHandler(r *mux.Router) {
	r.HandleFunc("/transactions", h.GetAllTransactions).Methods("GET")
	r.HandleFunc("/transactions", h.CreateTransaction).Methods("POST")
}

func (h *Handler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Inquiry] New Request All Transactions")

	transactions, err := h.service.GetAllTransactions()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not fetch transactions: %v", err), http.StatusInternalServerError)
		return
	}

	GenerateNewCreatedResponse(w, transactions, http.StatusOK)

	log.Printf("[Inquiry] Request All Transactions")
}

// func (h *Handler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil || id <= 0 {
// 		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
// 		return
// 	}

// 	transaction, err := h.service.GetTransactionByID(id)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Could not fetch transaction with ID %d: %v", id, err), http.StatusNotFound)
// 		return
// 	}

// 	response := NewSuccessResponse([]Transaction{*transaction})

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		log.Printf("Error encoding response: %v", err)
// 	}
// }

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	createdTransaction, err := h.service.CreateTransaction(t)
	if err != nil {
		// http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		GenerateNewErrorResponse(w, "Invalid Input data", "Bad Request", fmt.Sprintf("%v", r.Context().Value("requestID")), http.StatusBadRequest)
		return
	}

	GenerateNewCreatedResponse(w, createdTransaction, http.StatusOK)
}

// func (h *Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil || id <= 0 {
// 		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
// 		return
// 	}

// 	var t Transaction
// 	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
// 		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	updatedTransaction, err := h.service.UpdateTransaction(id, t)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Could not update transaction: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	response := NewSuccessResponse([]Transaction{updatedTransaction})

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		log.Printf("Error encoding response: %v", err)
// 	}
// }

// func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Query().Get("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil || id <= 0 {
// 		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = h.service.DeleteTransaction(id)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Could not delete transaction: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
// }
