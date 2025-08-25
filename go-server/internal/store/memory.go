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

// CreateWallet stores a new wallet
func (m *MemoryStore) CreateWallet(wallet *domain.Wallet) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wallets[wallet.ID] = wallet
}

// GetWallet retrieves a wallet by id
func (m *MemoryStore) GetWallet(id string) (*domain.Wallet, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	w, ok := m.wallets[id]
	return w, ok
}

// ListWallets returns all wallets
func (m *MemoryStore) ListWallets() []*domain.Wallet {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*domain.Wallet, 0, len(m.wallets))
	for _, w := range m.wallets {
		result = append(result, w)
	}
	return result
}

// CreateTransaction stores a new transaction
func (m *MemoryStore) CreateTransaction(tx *domain.Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transactions[tx.ID] = tx
}

// GetTransaction retrieves a transaction by id
func (m *MemoryStore) GetTransaction(id string) (*domain.Transaction, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tx, ok := m.transactions[id]
	return tx, ok
}

// ListTransactionsByWallet lists transactions for a wallet
func (m *MemoryStore) ListTransactionsByWallet(walletID string) []*domain.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := []*domain.Transaction{}
	for _, tx := range m.transactions {
		if tx.WalletID == walletID {
			result = append(result, tx)
		}
	}
	return result
}

// UpdateTransaction updates a transaction in place
func (m *MemoryStore) UpdateTransaction(tx *domain.Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transactions[tx.ID] = tx
}
