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
	if err := rc.Mx.Lock(); err != nil {
		// Lock is not granted.
		// either an internal error occured or there is an ongoing process hehe :D
		return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
			"code":    http.StatusTooManyRequests,
			"message": "too many request. try again later.",
		})
	}
	rc.Log.Info("lock granted", slog.String("ip request", c.IP()))
	// Release granted lock.
	defer func() {
		rc.Mx.Unlock()
		rc.Log.Info("lock released", slog.String("ip request", c.IP()))
	}()
	var req dto.BuyTicketRequest
	if err := c.Bind().Body(&req); err != nil {
		rc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	var ticket queries.Ticket
	err := queries.ExecTx(c.Context(), rc.Q.DB, func(q queries.Queries) error {
		ctx := c.Context()
		concert, err := q.Concert.GetConcert(ctx, req.ConcertID)
		if err != nil {
			return err
		}
		if concert.Limit == 0 {
			return errors.New("concert limit reached")
		}
		ticket, err = q.Ticket.CreateTicket(ctx, queries.CreateTicketQueryArgs{
			ConcertID:        req.ConcertID,
			TicketCategoryID: req.TicketCategoryID,
		})
		if err != nil {
			return err
		}
		concert.Limit--
		_, err = q.Concert.UpdateConcert(ctx, queries.UpdateConcertArgs{
			ID:    req.ConcertID,
			Limit: concert.Limit,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		rc.Log.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "failed to buy a ticket.",
			})
		}
		if errors.Is(err, errors.New("concert limit reached")) {
			return c.Status(http.StatusOK).JSON(fiber.Map{
				"code":    http.StatusOK,
				"message": "failed to buy a ticket. limit reached :(",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	rc.Log.Info("ticket created", "ticket", ticket)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "booking succeeded.",
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
