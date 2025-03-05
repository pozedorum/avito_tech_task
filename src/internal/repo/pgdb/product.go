package pgdb

import (
	"context"
	"fmt"

	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/pkg/postgres"
)

type ProductRepo struct {
	*postgres.Postgres
}

func NewProductRepo(pg *postgres.Postgres) *ProductRepo {
	return &ProductRepo{pg}
}

func (r *ProductRepo) CreateProduct(ctx context.Context, name string) (int, error) {
	sql, args, _ := r.Builder.
		Insert("products").
		Values(name).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ProductRepo.Deposit - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *ProductRepo) GetProductById(ctx context.Context, id int) (entity.Product, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("accounts").
		Where("id = ?", id).
		ToSql()
	var account entity.Product
	err := r.Pool.QueryRow(ctx, sql, args...).
		Scan(
			&account.Id,
			&account.Name,
		)
	if err != nil {
		return entity.Product{}, fmt.Errorf("ProductRepo.GetProductById - r.Pool.QueryRow: %v", err)
	}
	return account, nil
}

func (r *ProductRepo) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("accounts").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ProductRepo.GetAllProducts - r.Pool.QueryRow: %v", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err = rows.Scan(
			&product.Id,
			&product.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("ProductRepo.GetAllProducts - rows.Scan: %v", err)
		}
		products = append(products, product)
	}
	return products, nil
}
