package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

// ListWallets returns wallets (stub)
func (h *Handler) ListWallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]any{})
}

// CreateWallet creates a wallet (stub)
func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// GetWallet returns a wallet by id (stub)
func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// ListTransactions lists wallet transactions (stub)
func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode([]any{})
}

// CreateTransaction creates a transaction (stub)
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// GetTransaction returns a transaction (stub)
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// ApproveTransaction approves a transaction (stub)
func (h *Handler) ApproveTransaction(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}