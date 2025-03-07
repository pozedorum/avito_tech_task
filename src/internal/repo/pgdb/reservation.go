package pgdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/internal/repo/repoerrors"
	"github.com/pozedorum/user-balance-service/pkg/postgres"
)

type ReservationRepo struct {
	*postgres.Postgres
}

func NewReservationRepo(pg *postgres.Postgres) *ReservationRepo {
	return &ReservationRepo{pg}
}

func (r *ReservationRepo) CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - r.Pool.Begin: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sql, args, _ := r.Builder.Select("balance").
		From("accounts").
		Where("id = ?", reservation.AccountId).
		ToSql()
	var balance int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - r.Pool.QueryRow: %w", err)
	}

	if balance < reservation.Amount {
		return 0, repoerrors.ErrNotEnoughBalance
	}

	sql, args, _ = r.Builder.Update("accounts").
		Set("balance", sq.Expr("balance - ?", reservation.Amount)).
		Where("id = ?", reservation.AccountId).
		ToSql()
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.Exec(update accounts): %w", err)
	}

	sql, args, _ = r.Builder.Insert("reservations").
		Columns("account_id", "product_id", "order_id", "amount").
		Values(reservation.AccountId, reservation.ProductId, reservation.OrderId, reservation.Amount).
		Suffix("RETURNING id").
		ToSql()
	var id int
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.QueryRow: %w", err)
	}

	sql, args, _ = r.Builder.Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(reservation.AccountId, reservation.Amount, entity.OperationTypeReservation).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.Exec(insert operations): %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("ReservationRepo.CreateReservation - tx.Comit: %v", err)
	}

	return id, nil
}

func (r *ReservationRepo) GetReservationById(ctx context.Context, id int) (entity.Reservation, error) {
	sql, args, _ := r.Builder.Select("*").
		From("reservations").
		Where("id = ?", id).
		ToSql()
	var reservation entity.Reservation
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reservation.Id,
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
		&reservation.CreatedAt,
	)
	if err != nil {
		return entity.Reservation{}, fmt.Errorf("ReservationRepo.GetReservationById - r.Pool.QueryRow: %w", err)
	}
	return reservation, nil
}

func (r *ReservationRepo) RefundReservationById(ctx context.Context, id int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationById - r.Pool.Begin: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sql, args, _ := r.Builder.
		Delete("reservations").
		Where("id = ?", id).
		Suffix("RETURNING account_id, product_id, order_id, amount").
		ToSql()
	var reservation entity.Reservation
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
	)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationById - r.Pool.QueryRow: %w", err)
	}

	sql, args, _ = r.Builder.Update("accounts").
		Set("balance", sq.Expr("balance + ?", reservation.Amount)).
		Where("id = ?", reservation.AccountId).
		ToSql()
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationById - tx.Exec(update accounts): %w", err)
	}
	sql, args, _ = r.Builder.Insert("operations").
		Columns("account_id", "amount", "operation_type", "product_id", "order_id").
		Values(reservation.AccountId, reservation.Amount, entity.OperationTypeRefund, reservation.ProductId, reservation.OrderId).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationById - tx.Exec(insert operations): %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationById - tx.Comit: %v", err)
	}
	return nil
}

func (r *ReservationRepo) RefundReservationByOrderId(ctx context.Context, orderId int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - r.Pool.Begin: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sql, args, _ := r.Builder.
		Delete("reservations").
		Where("order_id = ?", orderId).
		Suffix("RETURNING account_id, product_id, order_id, amount").
		ToSql()
	var reservation entity.Reservation
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
	)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - r.Pool.QueryRow: %w", err)
	}

	sql, args, _ = r.Builder.Update("accounts").
		Set("balance", sq.Expr("balance + ?", reservation.Amount)).
		Where("id = ?", reservation.AccountId).
		ToSql()
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - tx.Exec(update accounts): %w", err)
	}
	sql, args, _ = r.Builder.Insert("operations").
		Columns("account_id", "amount", "operation_type", "product_id", "order_id").
		Values(reservation.AccountId, reservation.Amount, entity.OperationTypeRefund, reservation.ProductId, reservation.OrderId).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - tx.Exec(insert operations): %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - tx.Comit: %v", err)
	}
	return nil
}

func (r *ReservationRepo) RevenueReservationById(ctx context.Context, id int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - r.Pool.Begin: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sql, args, _ := r.Builder.
		Delete("reservations").
		Where("id = ?", id).
		Suffix("RETURNING account_id, product_id, order_id, amount").
		ToSql()
	var reservation entity.Reservation
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
	)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - r.Pool.QueryRow: %w", err)
	}

	sql, args, _ = r.Builder.Insert("operations").
		Columns("account_id", "amount", "operation_type", "product_id", "order_id").
		Values(reservation.AccountId, reservation.Amount, entity.OperationTypeRevenue, reservation.ProductId, reservation.OrderId).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - tx.Exec(insert operations): %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RevenueReservationById - tx.Comit: %v", err)
	}
	return nil
}

func (r *ReservationRepo) RevenueReservationByOrderId(ctx context.Context, orderId int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - r.Pool.Begin: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sql, args, _ := r.Builder.
		Delete("reservations").
		Where("order_id = ?", orderId).
		Suffix("RETURNING account_id, product_id, order_id, amount").
		ToSql()
	var reservation entity.Reservation
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.OrderId,
		&reservation.Amount,
	)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - r.Pool.QueryRow: %w", err)
	}

	sql, args, _ = r.Builder.Insert("operations").
		Columns("account_id", "amount", "operation_type", "product_id", "order_id").
		Values(reservation.AccountId, reservation.Amount, entity.OperationTypeRevenue, reservation.ProductId, reservation.OrderId).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - tx.Exec(insert operations): %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("ReservationRepo.RefundReservationByOrderId - tx.Comit: %v", err)
	}
	return nil
}
