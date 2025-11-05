package usecase

import (
	"errors"

	"github.com/av-plus-minus/homewok_go_05/ledger/internal/domain"
	"github.com/rs/zerolog/log"
)

type LedgerUseCase struct {
	budgetRepo      domain.BudgetRepository
	transactionRepo domain.TransactionRepository
}

func NewLedgerUseCase(budgetRepo domain.BudgetRepository, transactionRepo domain.TransactionRepository) *LedgerUseCase {
	return &LedgerUseCase{
		budgetRepo:      budgetRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *LedgerUseCase) AddTransaction(tx domain.Transaction) error {
	if err := tx.Validate(); err != nil {
		log.Error().Err(err).Str("category", tx.Category).Msg("transaction validation failed")
		return err
	}

	// Сохраняем оригинальное описание
	originalDescription := tx.Description

	// Если категории нет в бюджетах - используем категорию "другое"
	if _, exists := uc.budgetRepo.GetBudget(tx.Category); !exists {
		// Добавляем информацию о несуществующей категории в описание
		if originalDescription == "" {
			tx.Description = "(" + tx.Category + ")"
		} else {
			tx.Description = originalDescription + " (" + tx.Category + ")"
		}
		// Меняем категорию на "другое"
		tx.Category = "Другое"
		log.Warn().
			Msg("Расход сохранен в категории 'другое'.")
	}

	// ПРОВЕРКА БЮДЖЕТА ДОЛЖНА БЫТЬ ЗДЕСЬ - после возможного изменения категории
	if budget, exists := uc.budgetRepo.GetBudget(tx.Category); exists {
		currentSpent := uc.calculateCurrentSpending(tx.Category)

		if currentSpent+tx.Amount > budget.Limit {
			log.Warn().
				Str("category", tx.Category).
				Float64("limit", budget.Limit).
				Float64("current", currentSpent).
				Float64("new_amount", tx.Amount).
				Msg("budget exceeded")
			return errors.New("budget exceeded")
		}
	}

	tx.SetDefaultDate()

	if err := uc.transactionRepo.AddTransaction(tx); err != nil {
		log.Error().Err(err).Str("category", tx.Category).Msg("failed to add transaction")
		return err
	}

	log.Info().
		Str("category", tx.Category).
		Float64("amount", tx.Amount).
		Str("description", tx.Description).
		Msg("transaction added")
	return nil
}

func (uc *LedgerUseCase) ListTransactions() []domain.Transaction {
	return uc.transactionRepo.ListTransactions()
}

func (uc *LedgerUseCase) calculateCurrentSpending(category string) float64 {
	transactions := uc.transactionRepo.GetTransactionsByCategory(category)
	total := 0.0
	for _, tx := range transactions {
		total += tx.Amount
	}
	return total
}

// CheckValid - утилитарная функция для проверки валидации
func CheckValid(v interface{ Validate() error }) error {
	return v.Validate()
}
