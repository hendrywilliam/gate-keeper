package dto

import "github.com/hendrywilliam/gate-keeper/queries"

// Date time using Unix Epoch.

type CreateConcertRequest struct {
	Name     string `json:"name"`
	ArtistID int    `json:"artist_id"`
	Date     int    `json:"date"`
	VenueID  int    `json:"venue_id"`
	Limit    int    `json:"limit"`
}

type DeleteConcertRequest struct {
	ID queries.ConcertID `json:"id"`
}

type UpdateConcertRequest struct {
	ID        queries.ConcertID `json:"id"`
	Name      string            `json:"name"`
	ArtistID  int               `json:"artist_id"`
	VenueID   int               `json:"venue_id"`
	Date      int               `json:"date"`
	Limit     int               `json:"limit"`
	CreatedAt int               `json:"created_at"`
	UpdatedAt int               `json:"updated_at"`
}
