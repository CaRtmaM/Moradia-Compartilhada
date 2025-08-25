package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"wallet-backend/internal/domain"
	"wallet-backend/internal/store"
)

// Handler bundles dependencies and router
type Handler struct {
	r *chi.Mux
	s *store.MemoryStore
}

func New() *Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := &Handler{r: r, s: store.NewMemoryStore()}

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/wallets", func(r chi.Router) {
		r.Get("/", h.ListWallets)
		r.With(UserMiddleware).Post("/", h.CreateWallet)

		r.Route("/{walletID}", func(r chi.Router) {
			r.Get("/", h.GetWallet)
			r.Get("/transactions", h.ListTransactions)
			r.With(UserMiddleware).Post("/transactions", h.CreateTransaction)
		})
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/{txID}", h.GetTransaction)
		r.With(UserMiddleware).Post("/{txID}/approve", h.ApproveTransaction)
	})

	return h
}

func (h *Handler) Router() *chi.Mux {
	return h.r
}

// userCtxKey is used to store the user id in request context
type userCtxKey struct{}

// UserMiddleware ensures X-User-Id is present and adds it to context
func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-Id")
		if userID == "" {
			http.Error(w, "missing X-User-Id", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userCtxKey{}, userID)))
	})
}

// ListWallets returns wallets
func (h *Handler) ListWallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.s.ListWallets())
}

// CreateWallet creates a wallet
func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Name              string   `json:"name"`
		Owners            []string `json:"owners"`
		DailyLimitCents   int64    `json:"dailyLimitCents"`
		RequiredApprovals int      `json:"requiredApprovals"`
	}
	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if body.Name == "" || len(body.Owners) == 0 || body.RequiredApprovals < 1 {
		http.Error(w, "invalid fields", http.StatusBadRequest)
		return
	}
	if body.RequiredApprovals > len(body.Owners) {
		http.Error(w, "requiredApprovals cannot exceed owners", http.StatusBadRequest)
		return
	}

	newWallet := &domain.Wallet{
		ID:                uuid.NewString(),
		Name:              body.Name,
		Owners:            body.Owners,
		DailyLimitCents:   body.DailyLimitCents,
		RequiredApprovals: body.RequiredApprovals,
		CreatedAt:         time.Now().UTC(),
	}
	h.s.CreateWallet(newWallet)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newWallet)
}

// GetWallet returns a wallet by id
func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "walletID")
	wlt, ok := h.s.GetWallet(walletID)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(wlt)
}

// ListTransactions lists wallet transactions
func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "walletID")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.s.ListTransactionsByWallet(walletID))
}

// CreateTransaction creates a transaction
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "walletID")
	wlt, ok := h.s.GetWallet(walletID)
	if !ok {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}
	type req struct {
		AmountCents int64  `json:"amountCents"`
		Memo        string `json:"memo"`
	}
	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if body.AmountCents <= 0 {
		http.Error(w, "amount must be > 0", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(userCtxKey{}).(string)
	newTx := &domain.Transaction{
		ID:          uuid.NewString(),
		WalletID:    wlt.ID,
		AmountCents: body.AmountCents,
		Memo:        body.Memo,
		CreatedBy:   userID,
		Approvers:   []string{},
		Status:      "pending",
		CreatedAt:   time.Now().UTC(),
	}
	// Auto-approve if creator is an owner and amount under daily limit
	if containsString(wlt.Owners, userID) && body.AmountCents <= wlt.DailyLimitCents {
		newTx.Approvers = append(newTx.Approvers, userID)
		if len(newTx.Approvers) >= wlt.RequiredApprovals {
			newTx.Status = "approved"
		}
	}
	h.s.CreateTransaction(newTx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newTx)
}

// GetTransaction returns a transaction
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	txID := chi.URLParam(r, "txID")
	tx, ok := h.s.GetTransaction(txID)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tx)
}

// ApproveTransaction approves a transaction
func (h *Handler) ApproveTransaction(w http.ResponseWriter, r *http.Request) {
	txID := chi.URLParam(r, "txID")
	tx, ok := h.s.GetTransaction(txID)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	wlt, ok := h.s.GetWallet(tx.WalletID)
	if !ok {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}
	userID := r.Context().Value(userCtxKey{}).(string)
	if !containsString(wlt.Owners, userID) {
		http.Error(w, "not an owner", http.StatusForbidden)
		return
	}
	if containsString(tx.Approvers, userID) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tx)
		return
	}
	if tx.Status == "approved" || tx.Status == "executed" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tx)
		return
	}
	tx.Approvers = append(tx.Approvers, userID)
	if len(tx.Approvers) >= wlt.RequiredApprovals {
		tx.Status = "approved"
	}
	h.s.UpdateTransaction(tx)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tx)
}

func containsString(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}