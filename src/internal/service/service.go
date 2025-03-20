package service

import (
	"context"
	"time"

	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/internal/repo"
	"github.com/pozedorum/user-balance-service/internal/webapi"
	"github.com/pozedorum/user-balance-service/pkg/hasher"
)

// account.go
type AccountDepositInput struct {
	Id int
	Amount int
}

type AccountWithDrawInput struct {
	Id int
	Amount int
}

type AcocuntTransferInput struct {
	FromId int
	ToId int
	Amount int
}

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
	GetAccountById(ctx context.Context, userId int) (entity.Account, error)
	Deposit(ctx context.Context, input AccountDepositInput) error
	Withdraw(ctx context.Context, input AccountWithDrawInput) error
	Transfer(ctx context.Context, input AcocuntTransferInput) error
}


// auth.go
type AuthCreateUserInput struct {
	username string
	password string
}

type AuthGenerateTokenInput struct {
	username string
	password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParceToken(accessToken string) (int, error)
}

// operation.go
type OperationHistoryInput struct {
	AccountId int
	SortType  string
	Offset    int
	Limit     int
}

type OperationHistoryOutput struct {
	Amount      int       `json:"amount"`
	Operation   string    `json:"operation"`
	Time        time.Time `json:"time"`
	Product     string    `json:"product,omitempty"`
	Order       *int      `json:"order,omitempty"`
	Description string    `json:"description,omitempty"`
}

type Operation interface {
	OperationHistory(ctx context.Context, input OperationHistoryInput) ([]OperationHistoryOutput, error)
	MakeReportLink(ctx context.Context, month int, year int) (string, error)
	MakeReportFile(ctx context.Context, month int, year int) ([]byte, error)
}

// product.go

type Product interface {
	CreateProduct(ctx context.Context, name string) (int, error)
	GetProductById(ctx context.Context, id int) (entity.Product, error)
}


// reservation.go
type ReservationCreateInput struct {
	AccountId int
	ProductId int
	OrderId   int
	Amount    int
}

type Reservation interface {
	CreateReservation(ctx context.Context, input ReservationCreateInput) (int, error)
	RefundReservationById(ctx context.Context, id int) error
	RevenueReservationById(ctx context.Context, id int) error
}

type Services struct {
	Auth 		Auth
	Account 	Account
	Operation 	Operation
	Product 	Product
	Reservation Reservation
}

type ServicesDependenies struct {
	Repos *repo.Repositories
	GDrive webapi.GDrive
	Hasher hasher.PasswordHasher

	SignKey string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependenies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.User,deps.Hasher,deps.SignKey,deps.TokenTTL),
		Account: NewAccountService(deps.Repos.Account),
		Operation: NewOperationService(deps.Repos.Operation,deps.Repos.Product,deps.GDrive),
		Product: NewProductService(deps.Repos.Product),
		Reservation: NewReservationService(deps.Repos.Reservation),
	}
}