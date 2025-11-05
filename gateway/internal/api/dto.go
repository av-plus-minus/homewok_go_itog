package api

import (
	"time"

	"github.com/av-plus-minus/homewok_go_05/ledger/internal/domain"
)

type CreateTransactionRequest struct {
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
}

type TransactionResponse struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
}

type CreateBudgetRequest struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

type BudgetResponse struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

func (req CreateTransactionRequest) ToLedgerTx() domain.Transaction {
	return domain.Transaction{
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Date:        req.Date,
	}
}

func FromDomainTransaction(tx domain.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date,
	}
}

func (req CreateBudgetRequest) ToDomainBudget() domain.Budget {
	return domain.Budget{
		Category: req.Category,
		Limit:    req.Limit,
		Period:   req.Period,
	}
}

func FromDomainBudget(budget domain.Budget) BudgetResponse {
	return BudgetResponse{
		Category: budget.Category,
		Limit:    budget.Limit,
		Period:   budget.Period,
	}
}
