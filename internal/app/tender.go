package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"zadanie-6105/config"
	"zadanie-6105/internal/api"
	"zadanie-6105/internal/repo/postgres"
	"zadanie-6105/internal/service"
)

type App struct {
	api     *api.Api
	svc     *service.Service
	storage *postgres.Storage
	echo    *echo.Echo
}

func New(log *slog.Logger, cfg *config.Config) *App {
	app := &App{}

	app.echo = echo.New()

	db, err := postgres.ConnectPostgres(cfg)
	if err != nil {
		log.Error("failed to connect to PostgresSQL", err)
	}

	err = postgres.CreateTables(db)
	if err != nil {
		log.Error("failed to create table in DB", err)
	}

	app.storage = postgres.New(log, db)
	app.svc = service.New(log,
		app.storage, app.storage, app.storage,
		app.storage, app.storage, app.storage, app.storage, app.storage,
		app.storage)

	app.api = api.New(log,
		app.svc, app.svc, app.svc,
		app.svc, app.svc, app.svc, app.svc, app.svc,
	)

	app.echo.HTTPErrorHandler = customHTTPErrorHandler

	app.echo.GET("/api/ping", app.api.Ping)

	app.echo.GET("/api/tenders", app.api.Tenders)
	app.echo.POST("/api/tenders/new", app.api.CreateTender)
	app.echo.GET("/api/tenders/my", app.api.GetTenderByUser)
	app.echo.GET("/api/tenders/:tenderId/status", app.api.TenderStatus)
	app.echo.PUT("/api/tenders/:tenderId/status", app.api.ChangeTenderStatus)
	app.echo.PATCH("/api/tenders/:tenderId/edit", app.api.EditTender)
	app.echo.PUT("/api/tenders/:tenderId/rollback/:version", app.api.RollbackTender)

	app.echo.POST("/api/bids/new", app.api.CreateBid)
	app.echo.GET("/api/bids/my", app.api.GetBidsByUser)
	app.echo.GET("/api/bids/:tenderId/list", app.api.BidsForTender)
	app.echo.GET("/api/bids/:bidId/status", app.api.BidStatus)
	app.echo.PUT("/api/bids/:bidId/status", app.api.UpdateBidStatus)
	app.echo.PATCH("/api/bids/:bidId/edit", app.api.EditBid)
	app.echo.PUT("/api/bids/:bidId/submit_decision", app.api.SubmitDecision)
	app.echo.PUT("/api/bids/:bidId/feedback", app.api.Feedback)
	app.echo.PUT("/api/bids/:bidId/rollback/:version", app.api.RollbackBid)
	app.echo.GET("/api/bids/:tenderId/reviews", app.api.Reviews)

	return app
}

func (a *App) Run() error {
	fmt.Println("server running")

	err := a.echo.Start(":8080")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Println("stopping server..." + " op = app.Stop")

	if err := a.echo.Shutdown(ctx); err != nil {
		fmt.Println("failed to shutdown server")
		return err
	}
	return nil
}

// Custom error handler to change "message" to "reason".
func customHTTPErrorHandler(err error, ctx echo.Context) {
	response := map[string]interface{}{
		"reason": "Internal Server Error",
	}
	statusCode := http.StatusInternalServerError

	if httpError, ok := err.(*echo.HTTPError); ok {
		if httpError.Code != 0 {
			statusCode = httpError.Code
		}
		if httpError.Message != nil {
			if msg, ok := httpError.Message.(string); ok {
				response["reason"] = msg
			} else {
				response["reason"] = fmt.Sprintf("%v", httpError.Message)
			}
		} else {
			response["reason"] = httpError.Error()
		}
	} else {
		response["reason"] = err.Error()
	}

	ctx.JSON(statusCode, response)
}
