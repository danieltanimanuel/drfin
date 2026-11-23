package transactions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"errors"
	"drfin/app/accounts"
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

	// Enhancement: Use RESTful path param instead of query param (?id=)
	// This works better with HTMX hx-put / hx-delete
	r.HandleFunc("/transactions/{id}", h.UpdateTransaction).Methods("PUT", "PATCH")
	r.HandleFunc("/accounts/{id}/adjust-balance", h.AdjustAccountBalance).Methods("POST")
}

func (h *Handler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Inquiry] New Request All Transactions")

	// Add request-scoped timeout to prevent long-running DB queries
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var ownerID *int
	if v := r.URL.Query().Get("owner"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			http.Error(w, "invalid owner id", http.StatusBadRequest)
			return
		}
		ownerID = &id
	}

	transactions, err := h.service.GetAllTransactions(ctx, ownerID)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Could not fetch transactions: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	GenerateNewCreatedResponse(w, transactions, http.StatusOK)
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// Decode body first, then apply timeout only for business logic
	var t CreateTransactionDTO
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	createdTransaction, err := h.service.CreateTransaction(ctx, t)
	if err != nil {
		GenerateNewErrorResponse(
			w,
			"Invalid Input data",
			"Bad Request",
			err.Error(),
			fmt.Sprintf("%v", r.Context().Value("requestID")),
			http.StatusBadRequest,
		)
		return
	}

	GenerateNewCreatedResponse(w, createdTransaction, http.StatusOK)
}

func (h *Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Update] Transaction request")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Enhancement: Path params are safer & cleaner than query params
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		GenerateNewErrorResponse(
			w,
			"Invalid transaction ID",
			"Bad Request",
			fmt.Sprintf("invalid id: %v", idStr),
			fmt.Sprintf("%v", r.Context().Value("requestID")),
			http.StatusBadRequest,
		)
		return
	}

	var dto UpdateTransactionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		GenerateNewErrorResponse(
			w,
			"Invalid request body",
			"Bad Request",
			err.Error(),
			fmt.Sprintf("%v", r.Context().Value("requestID")),
			http.StatusBadRequest,
		)
		return
	}

	updated, err := h.service.UpdateTransaction(ctx, id, dto)
	if err != nil {
		GenerateNewErrorResponse(
			w,
			"Failed to update transaction",
			"Internal Server Error",
			err.Error(),
			fmt.Sprintf("%v", r.Context().Value("requestID")),
			http.StatusInternalServerError,
		)
		return
	}

	GenerateNewCreatedResponse(w, updated, http.StatusOK)
}

func (h *Handler) AdjustAccountBalance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		http.Error(w, "Invalid account id", http.StatusBadRequest)
		return
	}

	var req AdjustBalanceRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := h.service.AdjustAccountBalance(ctx, id, req)
	if err != nil {
		switch {
		case errors.Is(err, accounts.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, accounts.ErrAccountNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			log.Printf("Error adjusting balance: %v", err)
			http.Error(w, "Failed to adjust balance", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding adjustment result: %v", err)
	}
}