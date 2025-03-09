package service

import (
	"context"

	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/internal/repo"
)

type ReservationService struct {
	reservationRepo repo.Reservation
}

func NewReservationService(reservationRepo repo.Reservation) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo}
}

func (s *ReservationService) CreateReservation(ctx context.Context, input ReservationCreateInput) (int, error) {
	reservation := entity.Reservation{
		AccountId: input.AccountId,
		ProductId: input.ProductId,
		OrderId:   input.OrderId,
		Amount:    input.Amount,
	}

	id, err := s.reservationRepo.CreateReservation(ctx, reservation)
	if err != nil {
		return 0, ErrCannotCreateReservation
	}
	return id, nil
}

func (s *ReservationService) RefundReservationById(ctx context.Context, id int) error {
	return s.reservationRepo.RefundReservationById(ctx, id)
}

func (s *ReservationService) RevenueReservationById(ctx context.Context, id int) error {
	return s.reservationRepo.RefundReservationById(ctx, id)
}
