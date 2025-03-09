package service

import (
	"context"

	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/internal/repo"
)

type ProductService struct {
	productRepo repo.Product
}

func NewProductService(productRepo repo.Product) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string) (int, error) {
	return s.productRepo.CreateProduct(ctx, name)
}

func (s *ProductService) GetProductById(ctx context.Context, id int) (entity.Product, error) {
	return s.productRepo.GetProductById(ctx, id)
}
