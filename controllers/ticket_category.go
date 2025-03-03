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

type TicketCategoryController struct {
	Mx  *redsync.Mutex
	Q   *queries.Queries
	Log *slog.Logger
}

func NewTicketCategoryController(mx *redsync.Mutex, q *queries.Queries, log *slog.Logger) *TicketCategoryController {
	return &TicketCategoryController{
		Mx:  mx,
		Q:   q,
		Log: log,
	}
}

func (tc *TicketCategoryController) CreateTicketCategory(c fiber.Ctx) error {
	req := new(dto.CreateTicketCategoryRequest)
	if err := c.Bind().Body(req); err != nil {
		tc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	tcat, err := tc.Q.TicketCategory.CreateTicketCategory(c.Context(), queries.CreateTicketCategoryArgs{
		ConcertID:   req.ConcertID,
		Description: req.Description,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		tc.Log.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	tc.Log.Info("ticket category created", "ticket_category", tcat)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "concert updated.",
		"data":    tcat,
	})
}

func (tc *TicketCategoryController) UpdateTicketCategory(c fiber.Ctx) error {
	req := new(dto.UpdateTicketCategoryRequest)
	if err := c.Bind().Body(req); err != nil {
		tc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	tcat, err := tc.Q.TicketCategory.UpdateTicketCategory(c.Context(), queries.UpdateTicketCategoryArgs{
		ID:          req.ID,
		ConcertID:   req.ConcertID,
		Description: req.Description,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "no ticket category with the specified ID was found",
			})
		}
		tc.Log.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	tc.Log.Info("update ticket category succeeded", "ticket_category", tcat)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "ticket category updated.",
		"data":    tcat,
	})
}

func (tc *TicketCategoryController) DeleteTicketCategory(c fiber.Ctx) error {
	req := new(dto.DeleteTicketCategoryRequest)
	if err := c.Bind().Body(req); err != nil {
		tc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	tcat, err := tc.Q.TicketCategory.DeleteTicketCategory(c.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "no ticket category with the specified ID was found",
			})
		}
		tc.Log.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	tc.Log.Error("ticket category deleted", "ticket_category", tcat)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "ticket category deleted.",
		"data":    tcat,
	})
}
