package repo

import (
	"context"

	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/internal/repo/pgdb"
	"github.com/pozedorum/user-balance-service/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserById(ctx context.Context, id int) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
}

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
	GetAccountById(ctx context.Context, id int) (entity.Account, error)
	Deposit(ctx context.Context, id, amount int) error
	Withdraw(ctx context.Context, id int, amount int) error
	Transfer(ctx context.Context, from, to, amount int) error
}

type Product interface {
	CreateProduct(ctx context.Context, name string) (int, error)
	GetProductById(ctx context.Context, id int) (entity.Product, error)
	GetAllProducts(ctx context.Context) ([]entity.Product, error)
}

type Operation interface {
	GerAllRevenueOperationsGroupedByProduct(ctx context.Context, month, year int) ([]string, []int, error)
	OperationsPagination(ctx context.Context, accountId int, sortType string, offset int, limit int) ([]entity.Operation, []string, error)
}

type Reservation interface {
	CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error)
	GetReservationById(ctx context.Context, id int) (entity.Reservation, error)
	RefundReservationById(ctx context.Context, id int) error
	RefundReservationByOrderId(ctx context.Context, orderId int) error
	RevenueReservationById(ctx context.Context, id int) error
	RevenueReservationByOrderId(ctx context.Context, orderId int) error
}

type Repositories struct {
	User
	Account
	Product
	Operation
	Reservation
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:        pgdb.NewUserRepo(pg),
		Account:     pgdb.NewAccountRepo(pg),
		Product:     pgdb.NewProductRepo(pg),
		Operation:   pgdb.NewOperationRepo(pg),
		Reservation: pgdb.NewReservationRepo(pg),
	}
}
