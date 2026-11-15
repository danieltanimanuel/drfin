package accounts

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service  *Service
	renderer Renderer
}

type Renderer interface {
	Render(w http.ResponseWriter, templatePath string, templateName string,data interface{}) error
}

func NewHandler(service *Service, r Renderer) *Handler {
	return &Handler{service: service, renderer: r}
}

func (h *Handler) RegisterHandler(router *mux.Router) {
	router.HandleFunc("/accounts", h.GetAllAcounts).Methods("GET")
	router.HandleFunc(
        "/accounts/receivable",
        h.CreateReceivable,
    ).Methods("POST")
}



// GetTransactionAccountGroups handles GET requests to retrieve all transaction accountgroups
func (h *Handler) GetAllAcounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.GetAllAccounts()
	if err != nil {
		http.Error(w, "Failed to retrieve accounts", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(accounts)
    if err != nil {
        http.Error(w, "Failed to encode accounts", http.StatusInternalServerError)
        return
    }

	data := AccountPage{
        Accounts:     accounts,
        AccountsJSON: template.JS(jsonData), // mark safe for template
    }


	// Decide template based on HTMX header
	var templateName string
	if r.Header.Get("HX-Request") == "true" {
		templateName = "list_partial.html"
	} else {
		templateName = "list_full.html"
	}
	templatePath := "accounts/templates/"

	log.Println(accounts)
	err = h.renderer.Render(w, templatePath, templateName, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
    	http.Error(w, "Something error :(", http.StatusInternalServerError)
	}

	// Send the response wrapped in a success response format
	// response := NewSuccessResponse(accounts)

	// Send JSON response
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// if err := json.NewEncoder(w).Encode(response); err != nil {
	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// 	log.Printf("Error encoding response: %v", err)
	// }
}

func (h *Handler) CreateReceivable(
    w http.ResponseWriter,
    r *http.Request,
) {
    var payload struct {
        Name string `json:"name"`
    }

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if payload.Name == "" {
        http.Error(
            w,
            "Receivable name is required",
            http.StatusBadRequest,
        )
        return
    }

    account, err := h.service.CreateReceivable(ctx,payload.Name)
    if err != nil {
        log.Printf("Error creating receivable: %v", err)

        http.Error(
            w,
            "Failed to create receivable",
            http.StatusInternalServerError,
        )
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(account); err != nil {
        log.Printf("Error encoding account: %v", err)
    }
}
