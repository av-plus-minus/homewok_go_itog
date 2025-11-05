package domain

import "errors"

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

func (b Budget) Validate() error {
	if b.Category == "" {
		return errors.New("budget category must be non-empty")
	}
	if b.Limit <= 0 {
		return errors.New("budget limit must be > 0")
	}
	return nil
}

type BudgetRepository interface {
	SetBudget(budget Budget) error
	GetBudget(category string) (Budget, bool)
	GetAllBudgets() map[string]Budget
}

type TransactionRepository interface {
	AddTransaction(tx Transaction) error
	ListTransactions() []Transaction
	GetTransactionsByCategory(category string) []Transaction
}
