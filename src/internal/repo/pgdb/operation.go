package pgdb

import (
	"context"
	"fmt"

	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/pkg/postgres"
)

const (
	maxPaginationLimit     = 10
	defaultPaginationLimit = 10

	DataSortType   string = "date"
	AmountSortType string = "amount"
)

type OperationRepo struct {
	*postgres.Postgres
}

func NewOperationRepo(pg *postgres.Postgres) *OperationRepo {
	return &OperationRepo{pg}
}

func (r *OperationRepo) GerAllRevenueOperationsGroupedByProduct(ctx context.Context, month, year int) ([]string, []int, error) {
	sql, args, _ := r.Builder.Select("product.name", "sum(amount)").
		From("operations").
		InnerJoin("products on operations.product_id = products.id").
		Where("operation_type = ? and extract(month from operations.created_at) = ? and extract(year from operations.created_at) = ?", entity.OperationTypeRevenue, month, year).
		GroupBy("product.name").
		ToSql()
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("OperationRepo.GerAllRevenueOperationsGroupedByProduct - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var productNames []string
	var amounts []int

	for rows.Next() {
		var productName string
		var Amount int
		err = rows.Scan(&productName, &Amount)
		if err != nil {
			return nil, nil, fmt.Errorf("OperationRepo.GerAllRevenueOperationsGroupedByProduct - %w", err)
		}
		productNames = append(productNames, productName)
		amounts = append(amounts, Amount)
	}

	return productNames, amounts, nil
}

func (r *OperationRepo) OpertionsPagination(ctx context.Context, account int, sortType string, offset int, limit int) ([]entity.Operation, []string, error) {
}
