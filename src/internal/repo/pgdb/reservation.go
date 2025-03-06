package pgdb

import (
	// "context"
	// "errors"
	// "fmt"

	// "github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgconn"
	// "github.com/pozedorum/user-balance-service/internal/entity"
	// "github.com/pozedorum/user-balance-service/internal/repo/repoerrors"
	"github.com/pozedorum/user-balance-service/pkg/postgres"
	// log "github.com/sirupsen/logrus"
	// sq "github.com/Masterminds/squirrel"
)

type ReservationRepo struct {
	*postgres.Postgres
}

func NewReservationRepo(pg *postgres.Postgres) *ReservationRepo {
	return &ReservationRepo{pg}
}
