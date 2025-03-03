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

type ConcertController struct {
	Mx  *redsync.Mutex
	Q   *queries.Queries
	Log *slog.Logger
}

func NewConcertController(mx *redsync.Mutex, q *queries.Queries, log *slog.Logger) *ConcertController {
	return &ConcertController{
		Mx:  mx,
		Q:   q,
		Log: log,
	}
}

func (cc *ConcertController) CreateConcert(c fiber.Ctx) error {
	var req dto.CreateConcertRequest
	if err := c.Bind().Body(&req); err != nil {
		cc.Log.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process data",
		})
	}
	concert, err := cc.Q.Concert.CreateConcert(c.Context(), queries.CreateConcertQueryArgs{
		Name:     req.Name,
		ArtistID: req.ArtistID,
		VenueID:  req.VenueID,
		Date:     req.Date,
		Limit:    req.Limit,
	})
	if err != nil {
		cc.Log.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	cc.Log.Info("concert created", "concert", concert)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"code":    http.StatusCreated,
		"message": "concert created.",
		"data":    concert,
	})
}

func (cc *ConcertController) DeleteConcert(c fiber.Ctx) error {
	req := new(dto.DeleteConcertRequest)
	if err := c.Bind().Body(req); err != nil {
		slog.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process",
		})
	}
	concert, err := cc.Q.Concert.DeleteConcert(c.Context(), req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "no concert with the specified ID was found",
			})
		}
		slog.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	slog.Info("concert updated", "concert", concert)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "concert deleted.",
		"data":    concert,
	})
}

func (cc *ConcertController) UpdateConcert(c fiber.Ctx) error {
	req := new(dto.UpdateConcertRequest)
	if err := c.Bind().Body(req); err != nil {
		slog.Error(err.Error())
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "failed to process",
		})
	}
	concert, err := cc.Q.Concert.UpdateConcert(c.Context(), queries.UpdateConcertArgs{
		ID:       req.ID,
		Name:     req.Name,
		ArtistID: req.ArtistID,
		VenueID:  req.VenueID,
		Date:     req.Date,
		Limit:    req.Limit,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"message": "no concert with the specified ID was found",
			})
		}
		slog.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	}
	slog.Error("concert deleted", "concert", concert)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "concert updated.",
		"data":    concert,
	})
}
