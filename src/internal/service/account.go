package service

import (
	"context"

	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/internal/repo"
	"github.com/pozedorum/user-balance-service/internal/repo/repoerrors"
)

type AccountService struct {
	accountRepo repo.Account
}

func NewAccountService(accountRepo repo.Account) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(ctx context.Context) (int, error) {
	id, err := s.accountRepo.CreateAccount(ctx)
	if err != nil {
		if err == repoerrors.ErrAlreadyExists {
			return 0, ErrAccountAlreadyExists
		}
		return 0, ErrCannotCreateAccount
	}
	return id, nil
}

func (s *AccountService) GetAccountById(ctx context.Context, userId int) (entity.Account, error) {
	return s.accountRepo.GetAccountById(ctx, userId)
}

func (s *AccountService) Deposit(ctx context.Context, id int, amount int) error {
	return s.accountRepo.Deposit(ctx, id, amount)
}

func (s *AccountService) Withdraw(ctx context.Context, id int, amount int) error {
	return s.accountRepo.Withdraw(ctx, id, amount)
}

func (s *AccountService) Transfer(ctx context.Context, from int, to int, amount int) error {
	return s.accountRepo.Transfer(ctx, from, to, amount)
}
