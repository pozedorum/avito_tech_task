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

func (r *OperationRepo) OpertionsPagination(ctx context.Context, accountId int, sortType string, offset int, limit int) ([]entity.Operation, []string, error) {
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}
	if limit == 0 {
		limit = defaultPaginationLimit
	}

	var orderBySql string
	switch sortType {
	case "":
		orderBySql = "created_at DESC"
	case DataSortType:
		orderBySql = "created_at DESC"
	case AmountSortType:
		orderBySql = "amount DESC"
	default:
		return nil, nil, fmt.Errorf("OperationRepo.OperationsPagination - wrong sort type - %s", sortType)
	}

	sql, args, _ := r.Builder.Select("operations.id",
		"account_id",
		"amount",
		"operation_type",
		"created_at",
		"COALESCE((case when operations.product_id is null then null else products.name end), '') as product_name",
		"product_id",
		"order_id",
		"COALESCE(description, '')").
		From("operations").
		InnerJoin("products on products.product_id = operations.product_id or operations.product_id is null").
		Where("account_id = ?", accountId).
		OrderBy(orderBySql).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("OperationRepo.OperationsPagination - r.Pool.Query - %w", err)
	}
	defer rows.Close()

	var (
		operations   []entity.Operation
		productNames []string

		operation   entity.Operation
		productName string
	)
	for rows.Next() {
		err = rows.Scan(&operation.Id,
			&operation.AccountId,
			&operation.Amount,
			&operation.OperationType,
			&operation.CreatedAt,
			&productName,
			&operation.ProductId,
			&operation.OrderId,
			&operation.Description)
		if err != nil {
			return nil, nil, fmt.Errorf("OperationRepo.OperationsPagination - rows.Scan - %w", err)
		}
		operations = append(operations, operation)
		productNames = append(productNames, productName)
	}
	return operations, productNames, nil
}
