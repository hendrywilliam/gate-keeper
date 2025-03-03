package dto

// Date time using Unix Epoch.

type BuyTicketRequest struct {
	ConcertID        int `json:"concert_id"`
	TicketCategoryID int `json:"ticket_category"`
}

type GetTicketRequest struct {
	ID int `json:"id"`
}

type DeleteTicketRequest struct {
	ID int `json:"id"`
}
