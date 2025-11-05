package api

import (
	"encoding/json"
	"net/http"

	"github.com/av-plus-minus/homewok_go_05/ledger/internal/domain"
	"github.com/av-plus-minus/homewok_go_05/ledger/internal/usecase"
	"github.com/rs/zerolog/log"
)

type TransactionHandler struct {
	ledgerUC *usecase.LedgerUseCase
}

func NewTransactionHandler(ledgerUC *usecase.LedgerUseCase) *TransactionHandler {
	return &TransactionHandler{
		ledgerUC: ledgerUC,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest

	// Декодируем JSON тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("failed to decode transaction request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON format"})
		return
	}

	// Преобразуем в доменную модель
	tx := req.ToLedgerTx()

	// Валидируем транзакцию
	if err := tx.Validate(); err != nil {
		log.Warn().Err(err).Msg("transaction validation failed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Добавляем транзакцию через use case
	if err := h.ledgerUC.AddTransaction(tx); err != nil {
		log.Error().Err(err).Msg("failed to add transaction")
		w.Header().Set("Content-Type", "application/json")

		switch err.Error() {
		case "budget exceeded":
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "budget exceeded"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		}
		return
	}

	// Получаем добавленную транзакцию для ответа
	transactions := h.ledgerUC.ListTransactions()
	var createdTx domain.Transaction
	for _, t := range transactions {
		if t.Category == tx.Category && t.Amount == tx.Amount && t.Description == tx.Description {
			createdTx = t
			break
		}
	}

	// Формируем успешный ответ
	response := FromDomainTransaction(createdTx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	log.Info().
		Int("id", createdTx.ID).
		Str("category", createdTx.Category).
		Float64("amount", createdTx.Amount).
		Msg("transaction created successfully")
}
