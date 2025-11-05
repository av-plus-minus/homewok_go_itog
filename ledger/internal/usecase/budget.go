package usecase

import (
	"encoding/json"
	"io"

	"github.com/av-plus-minus/homewok_go_05/ledger/internal/domain"
	"github.com/rs/zerolog/log"
)

type BudgetUseCase struct {
	budgetRepo domain.BudgetRepository
}

func NewBudgetUseCase(budgetRepo domain.BudgetRepository) *BudgetUseCase {
	return &BudgetUseCase{budgetRepo: budgetRepo}
}

func (uc *BudgetUseCase) SetBudget(budget domain.Budget) error {
	if err := budget.Validate(); err != nil {
		log.Error().Err(err).Str("category", budget.Category).Msg("budget validation failed")
		return err
	}

	if err := uc.budgetRepo.SetBudget(budget); err != nil {
		log.Error().Err(err).Str("category", budget.Category).Msg("failed to set budget")
		return err
	}

	log.Info().
		Str("category", budget.Category).
		Float64("limit", budget.Limit).
		Str("period", budget.Period).
		Msg("budget set/updated")
	return nil
}

func (uc *BudgetUseCase) LoadBudgetsFromJSON(r io.Reader) error {
	var budgets []domain.Budget
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&budgets); err != nil {
		log.Error().Err(err).Msg("failed to parse budgets JSON")
		return err
	}

	for _, budget := range budgets {
		if err := uc.SetBudget(budget); err != nil {
			return err
		}
	}

	log.Info().Int("count", len(budgets)).Msg("budgets loaded from JSON")
	return nil
}

func (uc *BudgetUseCase) GetBudget(category string) (domain.Budget, bool) {
	return uc.budgetRepo.GetBudget(category)
}

func SeedBudgets(budgetUC *BudgetUseCase) {
	_ = budgetUC.SetBudget(domain.Budget{Category: "Тест_Категория_A", Limit: 9999, Period: "2025-09"})
	_ = budgetUC.SetBudget(domain.Budget{Category: "Тест_Категория_B", Limit: 999})
	_ = budgetUC.SetBudget(domain.Budget{Category: "Тест_Категория_C", Limit: 99})
	_ = budgetUC.SetBudget(domain.Budget{Category: "", Limit: 99999})
	_ = budgetUC.SetBudget(domain.Budget{Category: "Тест", Limit: -100})
}

func (uc *BudgetUseCase) ListBudgets() map[string]domain.Budget {
	return uc.budgetRepo.GetAllBudgets()
}
