package transaction

import (
	"context"
	"errors"
	"medium-rag/config"
)

type ITransactionService interface {
	getTotalTransactionPerMonth(ctx context.Context, query GetTotalTransactionReq) (GetTotalTransactionResp, error)
}

type ChatService struct {
	Env *config.EnvVariable
}

func NewTransactionService(env *config.EnvVariable) ITransactionService {

	return &ChatService{Env: env}
}

func (s *ChatService) getTotalTransactionPerMonth(ctx context.Context, query GetTotalTransactionReq) (GetTotalTransactionResp, error) {
	if query.Year == 2025 && query.Month == 2 {
		return GetTotalTransactionResp{
			TotalTransaction: 4000,
			Month:            query.Month,
			Year:             query.Year,
		}, nil
	}

	if query.Year == 2025 && query.Month == 1 {
		return GetTotalTransactionResp{
			TotalTransaction: 2000,
			Month:            query.Month,
			Year:             query.Year,
		}, nil
	}

	if query.Year == 2024 && query.Month == 12 {
		return GetTotalTransactionResp{
			TotalTransaction: 4000,
			Month:            query.Month,
			Year:             query.Year,
		}, nil
	}

	return GetTotalTransactionResp{}, errors.New("data not found")
}
