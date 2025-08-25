package store

import (
	"sync"

	"wallet-backend/internal/domain"
)

// MemoryStore is a naive in-memory store for prototyping
type MemoryStore struct {
	mu sync.RWMutex
	wallets map[string]*domain.Wallet
	transactions map[string]*domain.Transaction
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		wallets: make(map[string]*domain.Wallet),
		transactions: make(map[string]*domain.Transaction),
	}
}
