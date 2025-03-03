package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DbTx interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type Queries struct {
	DB     DbTx
	Ticket interface {
		CreateTicket(ctx context.Context, args CreateTicketQueryArgs) (Ticket, error)
		DeleteTicket(ctx context.Context, id TicketID) (Ticket, error)
		GetTicket(ctx context.Context, id TicketID) (Ticket, error)
	}
	Concert interface {
		CreateConcert(ctx context.Context, args CreateConcertQueryArgs) (Concert, error)
		DeleteConcert(ctx context.Context, id ConcertID) (Concert, error)
		UpdateConcert(ctx context.Context, args UpdateConcertArgs) (Concert, error)
	}
	TicketCategory interface {
		UpdateTicketCategory(ctx context.Context, args UpdateTicketCategoryArgs) (TicketCategory, error)
		DeleteTicketCategory(ctx context.Context, id TicketCategoryID) (TicketCategory, error)
		CreateTicketCategory(ctx context.Context, args CreateTicketCategoryArgs) (TicketCategory, error)
	}
}

func NewQueries(db DbTx) Queries {
	return Queries{
		DB:             db,
		Concert:        &ConcertQueryImpl{DB: db},
		TicketCategory: &TicketCategoryQueryImpl{DB: db},
		Ticket:         &TicketQueryImpl{DB: db},
	}
}

func ExecTx(ctx context.Context, db DbTx, fn func(Queries) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	n := NewQueries(tx)
	if err = fn(n); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
