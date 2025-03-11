package service

import (
	"context"

	"github.com/pozedorum/user-balance-service/internal/webapi"

	"github.com/pozedorum/user-balance-service/internal/repo"
)

type OperationService struct {
	operationRepo repo.Operation
	productRepo   repo.Product
	gDrive        webapi.GDrive
}

func NewOperationService(operationRepo repo.Operation, productRepo repo.Product, gDrive webapi.GDrive) *OperationService {
	return &OperationService{
		operationRepo: operationRepo,
		productRepo:   productRepo,
		gDrive:        gDrive,
	}
}

func (s *OperationService) OperationHistory(ctx context.Context, input OperationHistoryInput) ([]OperationHistoryOutput, error) {
	operations, productNames, err := s.operationRepo.OperationsPagination(ctx, input.AccountId, input.SortType, input.Offset, input.Limit)
	if err != nil {
		return nil, err
	}

	output := make([]OperationHistoryOutput, 0, len(operations))

	for i, operation := range operations {
		output = append(output, OperationHistoryOutput{
			Amount:    operation.Amount,
			Operation: operation.OperationType,
			Time:      operation.CreatedAt,
			Product: productNames[i],
			Order: operation.OrderId,
			Description: operation.Description,
		})
	}
	return output, nil
}
