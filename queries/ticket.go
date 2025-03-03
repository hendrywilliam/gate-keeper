package queries

import (
	"context"
	"log/slog"
	"time"
)

type TicketID = int

type Ticket struct {
	ID               TicketID `json:"id"`
	SerialNumber     string   `json:"serial_number"`
	ConcertID        int      `json:"concert_id"`
	TicketCategoryID int      `json:"ticket_category"`
	CreatedAt        int      `json:"created_at,omitempty"`
	UpdatedAt        int      `json:"updated_at,omitempty"`
}

func (t Ticket) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("ticket_id", t.ID),
		slog.String("serial_number", t.SerialNumber),
	)
}

type TicketQueryImpl struct {
	DB DbTx
}

type CreateTicketQueryArgs struct {
	ConcertID        int
	TicketCategoryID int
}

func (tq *TicketQueryImpl) CreateTicket(ctx context.Context, args CreateTicketQueryArgs) (Ticket, error) {
	row := tq.DB.QueryRow(ctx, `
		INSERT INTO ticket (
			concert_id,
			ticket_category_id,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING id, serial_number, concert_id, ticket_category_id;
	`, args.ConcertID, args.TicketCategoryID, time.Now().Unix(), time.Now().Unix())
	var t Ticket
	err := row.Scan(
		&t.ID,
		&t.SerialNumber,
		&t.ConcertID,
		&t.TicketCategoryID,
	)
	return t, err
}

func (tq *TicketQueryImpl) GetTicket(ctx context.Context, id TicketID) (Ticket, error) {
	row := tq.DB.QueryRow(ctx, `
		SELECT
			serial_number,
			concert_id,
			ticket_category_id
		FROM ticket
		WHERE id = $1;	
	`, id)
	var t Ticket
	err := row.Scan(
		&t.SerialNumber,
		&t.ConcertID,
		&t.TicketCategoryID,
	)
	return t, err
}

func (tq *TicketQueryImpl) DeleteTicket(ctx context.Context, id TicketID) (Ticket, error) {
	row := tq.DB.QueryRow(ctx, `
		DELETE FROM ticket
		WHERE id = $1
		RETURNING id, serial_number;
	`, id)
	var t Ticket
	err := row.Scan(
		&t.ID,
		&t.SerialNumber,
	)
	return t, err
}
