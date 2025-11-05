package store

import (
	"sync"

	"github.com/av-plus-minus/homewok_go_05/ledger/internal/domain"
)

type InMemoryStore struct {
	mu           sync.RWMutex
	transactions []domain.Transaction
	budgets      map[string]domain.Budget
	nextID       int
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		transactions: make([]domain.Transaction, 0, 32),
		budgets:      make(map[string]domain.Budget, 8),
		nextID:       1,
	}
}

// Реализация BudgetRepository
func (s *InMemoryStore) SetBudget(budget domain.Budget) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.budgets[budget.Category] = budget
	return nil
}

func (s *InMemoryStore) GetBudget(category string) (domain.Budget, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	budget, exists := s.budgets[category]
	return budget, exists
}

func (s *InMemoryStore) GetAllBudgets() map[string]domain.Budget {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]domain.Budget)
	for k, v := range s.budgets {
		result[k] = v
	}
	return result
}

// Реализация TransactionRepository
func (s *InMemoryStore) AddTransaction(tx domain.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx.ID = s.nextID
	s.nextID++
	s.transactions = append(s.transactions, tx)
	return nil
}

func (s *InMemoryStore) ListTransactions() []domain.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Transaction, len(s.transactions))
	copy(result, s.transactions)
	return result
}

func (s *InMemoryStore) GetTransactionsByCategory(category string) []domain.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Transaction
	for _, tx := range s.transactions {
		if tx.Category == category {
			result = append(result, tx)
		}
	}
	return result
}
