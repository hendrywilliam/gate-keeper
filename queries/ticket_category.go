package queries

import (
	"context"
	"log/slog"
	"time"
)

type TicketCategoryID = int

type TicketCategory struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ConcertID   int     `json:"concert_id"`
	StartDate   int     `json:"start_date"`
	EndDate     int     `json:"end_date"`
	CreatedAt   int     `json:"created_at"`
	UpdatedAt   int     `json:"updated_at"`
}

func (tc TicketCategory) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", tc.ID),
		slog.String("description", tc.Description),
		slog.Float64("price", tc.Price),
		slog.Int("concert_id", tc.ConcertID),
	)
}

type TicketCategoryQueryImpl struct {
	DB DbTx
}

type CreateTicketCategoryArgs struct {
	ConcertID   int
	Description string
	Price       float64
	StartDate   int
	EndDate     int
}

func (tc *TicketCategoryQueryImpl) CreateTicketCategory(ctx context.Context, args CreateTicketCategoryArgs) (TicketCategory, error) {
	row := tc.DB.QueryRow(ctx, `
		INSERT INTO ticket_category (
			concert_id,
			description,
			price,
			start_date,
			end_date,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7 
		) RETURNING id, concert_id, description, price;
	`, args.ConcertID, args.Description, args.Price, args.StartDate, args.EndDate, time.Now().Unix(), time.Now().Unix())
	var tcat TicketCategory
	err := row.Scan(
		&tcat.ID,
		&tcat.ConcertID,
		&tcat.Description,
		&tcat.Price,
	)
	return tcat, err
}

type UpdateTicketCategoryArgs struct {
	ID          TicketCategoryID
	ConcertID   int
	Description string
	Price       float64
	StartDate   int
	EndDate     int
}

func (tc *TicketCategoryQueryImpl) UpdateTicketCategory(ctx context.Context, args UpdateTicketCategoryArgs) (TicketCategory, error) {
	row := tc.DB.QueryRow(ctx, `
		UPDATE ticket_category
		SET concert_id = $1,
			description = $2,
			price = $3,
			start_date = $4,
			end_date = $5
		WHERE id = $6
		RETURNING id, concert_id, description, price, start_date, end_date;
	`, args.ConcertID, args.Description, args.Price, args.StartDate, args.EndDate, args.ID)
	var tcat TicketCategory
	err := row.Scan(
		&tcat.ID,
		&tcat.ConcertID,
		&tcat.Description,
		&tcat.Price,
		&tcat.StartDate,
		&tcat.EndDate,
	)
	return tcat, err
}

func (tc *TicketCategoryQueryImpl) DeleteTicketCategory(ctx context.Context, id TicketCategoryID) (TicketCategory, error) {
	row := tc.DB.QueryRow(ctx, `
		DELETE FROM ticket_category
		WHERE id = $1
		RETURNING id, description, price;
	`)
	var tcat TicketCategory
	err := row.Scan(
		&tcat.ID,
		&tcat.Description,
		&tcat.Price,
	)
	return tcat, err
}
