package transactionstype

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

type TransactionsType struct {
	transactionTypeID       int
	name                    string
	parentTransactionTypeID sql.NullInt16
}

type Response struct {
	Status           string             `json: status`
	Message          string             `json: message`
	TransactionsType []TransactionsType `json: data`
}

func NewResponse(Status string, Message string, TransactionsType []TransactionsType) *Response {
	return &Response{
		Status:           Status,
		Message:          Message,
		TransactionsType: TransactionsType,
	}
}

func (h *Handler) RegisterHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /transactions/type", h.RetrieveTransactionType)
	mux.HandleFunc("POST /transactions/type", h.InsertTransactionType)
	mux.HandleFunc("GET /transactions/type/{transactionTypeID}", h.RetrieveTransactionTypeByID)
}

func (h *Handler) RetrieveTransactionType(w http.ResponseWriter, r *http.Request) {
	selectstmt := `SELECT transaction_type_id, name, parent_transaction_type_id FROM transactions_type`
	rows, err := h.db.Query(selectstmt)
	if err != nil {
		log.Panic(err)
		return
	}
	defer rows.Close()
	var transactionTypes []TransactionsType
	for rows.Next() {
		var transactionType TransactionsType
		if err = rows.Scan(&transactionType.transactionTypeID, &transactionType.name, &transactionType.parentTransactionTypeID); err != nil {
			log.Panic(err)
		}
		transactionTypes = append(transactionTypes, transactionType)
	}
	response := NewResponse("200", "Success", transactionTypes)
	log.Println(response)
	var jsonData []byte
	jsonData, err = json.Marshal(response)
	if err != nil {
		log.Panic(err)
		return
	}
	var jsonString = string(jsonData)
	fmt.Fprintf(w, "%+v", jsonString)
}

func (h *Handler) RetrieveTransactionTypeByID(w http.ResponseWriter, r *http.Request) {
	transactionID := r.PathValue("transactionTypeID")
	selectstmt := `SELECT transaction_type_id, name, parent_transaction_type_id FROM transactions_type WHERE transaction_type_id = $1`
	var transactionType TransactionsType
	err := h.db.QueryRow(selectstmt, transactionID).Scan(&transactionType.transactionTypeID, &transactionType.name, &transactionType.parentTransactionTypeID)
	jsonTransactionType := new(bytes.Buffer)
	json.NewEncoder(jsonTransactionType).Encode(transactionType)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			fmt.Fprintln(w, jsonTransactionType)
			return
		}
		log.Panic(err)
	}
	fmt.Fprintln(w, "Transaction Type : ", jsonTransactionType)
}

func (h *Handler) InsertTransactionType(w http.ResponseWriter, r *http.Request) {
	var transactionType TransactionsType
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&transactionType)
	if err != nil {
		log.Panic(err)
	}
	insertstmt := `INSERT INTO transactions_type (name, parent_transaction_type_id) values ($1,$2) RETURNING transaction_type_id`
	transactionTypeID := 0
	err = h.db.QueryRow(insertstmt, transactionType.name, transactionType.parentTransactionTypeID).Scan(&transactionTypeID)
	if err != nil {
		log.Panic(err)
	}
	transactionType.transactionTypeID = transactionTypeID
	fmt.Fprintf(w, "Transaction Type Inserted Successfully %+v", transactionType)
}
