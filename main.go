package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	cfg "github.com/hendrywilliam/gate-keeper/config"
	"github.com/hendrywilliam/gate-keeper/controllers"
	"github.com/hendrywilliam/gate-keeper/queries"
	"github.com/hendrywilliam/gate-keeper/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	var logHandler slog.Handler
	logOpts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	if os.Getenv("APP_ENV") != "production" {
		logHandler = utils.NewCustomHandler(os.Stdout, utils.CustomHandlerOpts{
			SlogOpts: logOpts,
		})
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &logOpts)
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	app := fiber.New()
	mutex := cfg.NewMutex()
	db, err := cfg.NewPg(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to open pg")
		os.Exit(1)
	}
	// Move to heap (as long live object).
	allQs := queries.NewQueries(db)

	concertCtrl := controllers.NewConcertController(mutex, &allQs, logger)
	ticketCtrl := controllers.NewTicketController(mutex, &allQs, logger)
	tcatCtrl := controllers.NewTicketCategoryController(mutex, &allQs, logger)

	app.Post("/concert", concertCtrl.CreateConcert)
	app.Delete("/concert", concertCtrl.DeleteConcert)
	app.Put("/concert", concertCtrl.UpdateConcert)

	app.Post("/ticket", ticketCtrl.BuyTicket)
	app.Delete("/ticket", ticketCtrl.CancelTicket)
	app.Get("/ticket", ticketCtrl.GetTicket)

	app.Post("/ticket-category", tcatCtrl.CreateTicketCategory)
	app.Put("/ticket-category", tcatCtrl.UpdateTicketCategory)
	app.Delete("/ticket-category", tcatCtrl.DeleteTicketCategory)

	app.Listen(":8080")
}
