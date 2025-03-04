package queries

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
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

func (cq *ConcertQueryImpl) GetConcert(ctx context.Context, ID ConcertID) (Concert, error) {
	row := cq.DB.QueryRow(ctx, `
		SELECT
			name,
			artist_id,
			venue_id,
			date,
			"limit"
		FROM concert
		WHERE id = $1;
	`, ID)
	var t Concert
	err := row.Scan(
		&t.Name,
		&t.ArtistID,
		&t.VenueID,
		&t.Date,
		&t.Limit,
	)
	return t, err
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
	baseSql := &bytes.Buffer{}
	var setClauses []string
	var arguments []interface{}
	paramIndex := 1
	baseSql.WriteString("UPDATE concert ")
	// Empty field has its own default value.
	if args.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%v", paramIndex))
		arguments = append(arguments, args.Name)
		paramIndex++
	}
	if args.ArtistID > 0 {
		setClauses = append(setClauses, fmt.Sprintf("artist_id = $%v", paramIndex))
		arguments = append(arguments, args.ArtistID)
		paramIndex++
	}
	if args.VenueID > 0 {
		setClauses = append(setClauses, fmt.Sprintf("artist_id = $%v", paramIndex))
		arguments = append(arguments, args.VenueID)
		paramIndex++
	}
	// Date time is using epoch.
	if args.Date > 0 {
		setClauses = append(setClauses, fmt.Sprintf("date = $%v", paramIndex))
		arguments = append(arguments, args.Date)
		paramIndex++
	}
	if args.Limit >= 0 {
		setClauses = append(setClauses, fmt.Sprintf(`"limit" = $%v`, paramIndex))
		arguments = append(arguments, args.Limit)
		paramIndex++
	}
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%v", paramIndex))
	arguments = append(arguments, time.Now().Unix())
	paramIndex++
	baseSql.WriteString(fmt.Sprintf("SET %s ", strings.Join(setClauses, ", ")))
	baseSql.WriteString(fmt.Sprintf("WHERE id = $%v RETURNING id, name;", paramIndex))
	arguments = append(arguments, args.ID)
	fmt.Println(baseSql.String(), arguments)
	row := cq.DB.QueryRow(ctx, baseSql.String(), arguments...)
	var c Concert
	err := row.Scan(
		&c.ID,
		&c.Name,
	)
	return c, err
}
