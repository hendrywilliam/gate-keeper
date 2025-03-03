package queries

import (
	"context"
	"log/slog"
	"time"
)

type ConcertID = int

type Concert struct {
	ID        ConcertID `json:"id,omitempty"`
	Name      string    `json:"concert,omitempty"`
	ArtistID  int       `json:"artist_id,omitempty"`
	Date      int       `json:"date,omitempty"`
	VenueID   int       `json:"venue_id,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	CreatedAt int       `json:"created_at,omitempty"`
	UpdatedAt int       `json:"updated_at,omitempty"`
}

func (c Concert) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", c.ID),
		slog.String("name", c.Name),
		slog.Int("date", c.Date),
		slog.Int("limit", c.Limit),
	)
}

type ConcertQueryImpl struct {
	DB DbTx
}

type CreateConcertQueryArgs struct {
	Name     string
	ArtistID int
	VenueID  int
	Date     int
	Limit    int
}

func (cq *ConcertQueryImpl) CreateConcert(ctx context.Context, args CreateConcertQueryArgs) (Concert, error) {
	row := cq.DB.QueryRow(ctx, `
		INSERT INTO concert (
			name,
			artist_id,
			venue_id,
			date,
			"limit",
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id, name, date, "limit";
	`, args.Name, args.ArtistID, args.VenueID, args.Date, args.Limit, time.Now().Unix(), time.Now().Unix())
	var c Concert
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Date,
		&c.Limit,
	)
	return c, err
}

func (cq *ConcertQueryImpl) DeleteConcert(ctx context.Context, id ConcertID) (Concert, error) {
	row := cq.DB.QueryRow(ctx, `
		DELETE FROM concert
		WHERE id = $1
		RETURNING id, name;
	`, id)
	var c Concert
	err := row.Scan(
		&c.ID,
		&c.Name,
	)
	return c, err
}

type UpdateConcertArgs struct {
	ID        ConcertID
	Name      string
	ArtistID  int
	VenueID   int
	Date      int
	Limit     int
	CreatedAt int
	UpdatedAt int
}

func (cq *ConcertQueryImpl) UpdateConcert(ctx context.Context, args UpdateConcertArgs) (Concert, error) {
	row := cq.DB.QueryRow(ctx, `
		UPDATE concert
		SET name = $1,
			artist_id = $2,
			venue_id = $3,
			date = $4,
			"limit" = $5,
			updated_at = $6
		WHERE id = $7
		RETURNING id, name;
	`, args.Name, args.ArtistID, args.VenueID, args.Date, args.Limit, time.Now().Unix(), args.ID)
	var c Concert
	err := row.Scan(
		&c.ID,
		&c.Name,
	)
	return c, err
}
