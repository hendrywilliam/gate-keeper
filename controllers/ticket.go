package controllers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-redsync/redsync/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/hendrywilliam/gate-keeper/dto"
	"github.com/hendrywilliam/gate-keeper/queries"
)

type TicketController struct {
	Mx  *redsync.Mutex
	Q   *queries.Queries
	Log *slog.Logger
}

func NewTicketController(mx *redsync.Mutex, q *queries.Queries, log *slog.Logger) *TicketController {
	return &TicketController{
		Mx:  mx,
		Q:   q,
		Log: log,
	}
}

func (rc *TicketController) BuyTicket(c fiber.Ctx) error {
	var req dto.BuyTicketRequest
	if err := c.Bind().Body(&req); err != nil {
		rc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	ticket, err := rc.Q.Ticket.CreateTicket(c.Context(), queries.CreateTicketQueryArgs{
		ConcertID:        req.ConcertID,
		TicketCategoryID: req.TicketCategoryID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "failed to buy a ticket.",
			})
		}
		rc.Log.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	rc.Log.Error("ticket created", "ticket", ticket)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "concert deleted.",
		"data":    ticket,
	})
}

func (rc *TicketController) GetTicket(c fiber.Ctx) error {
	var req dto.GetTicketRequest
	if err := c.Bind().Body(&req); err != nil {
		rc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	ticket, err := rc.Q.Ticket.GetTicket(c.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "no ticket with the specified ID was found",
			})
		}
		rc.Log.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	rc.Log.Error("ticket obtained", "ticket", ticket)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "ticket obtained.",
		"data":    ticket,
	})
}

func (rc *TicketController) CancelTicket(c fiber.Ctx) error {
	var req dto.GetTicketRequest
	if err := c.Bind().Body(&req); err != nil {
		rc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	ticket, err := rc.Q.Ticket.DeleteTicket(c.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "no ticket with the specified ID was found",
			})
		}
		rc.Log.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	rc.Log.Error("ticket deleted", "ticket", ticket)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "ticket canceled.",
		"data":    ticket,
	})
}
