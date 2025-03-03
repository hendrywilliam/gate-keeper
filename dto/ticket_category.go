package dto

type CreateTicketCategoryRequest struct {
	ConcertID   int     `json:"concert_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StartDate   int     `json:"start_date"`
	EndDate     int     `json:"end_date"`
}

type UpdateTicketCategoryRequest struct {
	ID          int     `json:"id"`
	ConcertID   int     `json:"concert_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StartDate   int     `json:"start_date"`
	EndDate     int     `json:"end_date"`
}

type DeleteTicketCategoryRequest struct {
	ID int `json:"id"`
}
