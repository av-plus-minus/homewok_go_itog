package main

import (
	"net/http"
	"os"

	"gateway/internal/api"

	"github.com/av-plus-minus/homewok_go_05/ledger/internal/store"
	"github.com/av-plus-minus/homewok_go_05/ledger/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = "2006-01-02"
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02",
	})

	// Инициализация хранилища и use cases
	storage := store.NewInMemoryStore()
	ledgerUC := usecase.NewLedgerUseCase(storage, storage)

	// Инициализация обработчиков
	transactionHandler := api.NewTransactionHandler(ledgerUC)

	// Создаем роутер chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Базовые маршруты
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	// API маршруты с префиксом /api
	r.Route("/api", func(r chi.Router) {
		// Транзакции
		r.Route("/transactions", func(r chi.Router) {
			r.Post("/", transactionHandler.CreateTransaction) // Создание транзакции
			r.Get("/", nil)                                   // Список всех транзакций
			r.Get("/{id}", nil)                               // Получение транзакции по ID
			r.Delete("/{id}", nil)                            // Удаление транзакции
		})

		// Бюджеты
		r.Route("/budgets", func(r chi.Router) {
			r.Post("/", nil)             // Создание/обновление бюджета
			r.Get("/", nil)              // Список всех бюджетов
			r.Get("/{category}", nil)    // Получение бюджета по категории
			r.Delete("/{category}", nil) // Удаление бюджета
		})

		// Отчеты
		r.Route("/reports", func(r chi.Router) {
			r.Get("/spending", nil) // Отчет по тратам
			r.Get("/budgets", nil)  // Отчет по бюджетам
		})
	})

	// Статус страница
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok", "service": "gateway"}`))
	})

	addr := ":8080"
	log.Info().Str("port", addr).Msg("Gateway API server starting")
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Error().Err(err).Msg("gateway server error")
	}
}
