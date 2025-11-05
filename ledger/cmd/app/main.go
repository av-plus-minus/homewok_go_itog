package main

import (
	"bufio"
	"os"
	"time"

	"github.com/av-plus-minus/homewok_go_05/ledger/internal/domain"
	"github.com/av-plus-minus/homewok_go_05/ledger/internal/store"
	"github.com/av-plus-minus/homewok_go_05/ledger/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Настройка zerolog
	zerolog.TimeFieldFormat = "2006-01-02" //zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02",
	})

	log.Info().Msg("Ledger service started")

	// Инициализация хранилища и use cases
	storage := store.NewInMemoryStore()
	budgetUC := usecase.NewBudgetUseCase(storage)
	ledgerUC := usecase.NewLedgerUseCase(storage, storage)

	// Загрузка бюджетов
	if file, err := os.Open("data.json"); err == nil {
		defer file.Close()
		log.Info().Msg("Loading budgets from data.json")
		if err := budgetUC.LoadBudgetsFromJSON(bufio.NewReader(file)); err != nil {
			log.Error().Err(err).Msg("failed to load budgets, using defaults")
			usecase.SeedBudgets(budgetUC)
		}
	} else {
		log.Warn().Msg("data.json not found, using data defaults")
		usecase.SeedBudgets(budgetUC)
	}

	// Демонстрация работы
	budget := domain.Budget{Category: "Тест", Limit: 300, Period: "2025-01-13"}
	if err := usecase.CheckValid(budget); err != nil {
		log.Error().Err(err).Msg("CheckValid budget failed")
	} else {
		log.Info().Str("category", budget.Category).Msg("CheckValid budget OK")
	}

	tx := domain.Transaction{Amount: 11, Category: "Тест", Description: "Тестовое описание", Date: time.Now()}
	if err := usecase.CheckValid(tx); err != nil {
		log.Error().Err(err).Msg("CheckValid transaction failed")
	} else {
		log.Info().Str("category", tx.Category).Msg("CheckValid transaction OK")
	}

	// Добавление транзакций
	if err := ledgerUC.AddTransaction(tx); err != nil {
		log.Error().Err(err).Msg("AddTransaction A failed")
	}

	if err := ledgerUC.AddTransaction(domain.Transaction{
		Amount:      1_000_000,
		Category:    "тест_3",
		Description: "Подозрительный расход",
	}); err != nil {
		log.Warn().Err(err).Msg("AddTransaction B expected error")
	}

	if err := ledgerUC.AddTransaction(domain.Transaction{
		Amount:      400,
		Category:    "транспорт",
		Description: "Такси",
	}); err != nil {
		log.Error().Err(err).Msg("AddTransaction C failed")
	}

	// Вывод транзакций
	for _, t := range ledgerUC.ListTransactions() {
		log.Info().
			Int("id", t.ID).
			Str("category", t.Category).
			Float64("amount", t.Amount).
			Str("desc", t.Description).
			Time("date", t.Date).
			Msg("transaction listed")
	}

	// Вывод списка доступных категорий бюджетов
	log.Info().Msg("=== Доступные категории бюджетов ===")
	budgets := budgetUC.ListBudgets()
	if len(budgets) == 0 {
		log.Info().Msg("Нет установленных бюджетов")
	} else {
		for category, budget := range budgets {
			log.Info().
				Str("категория", category).
				Float64("лимит", budget.Limit).
				Str("период", budget.Period).
				Msg("бюджет")
		}
	}
}
