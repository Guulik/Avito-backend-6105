package api

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type Api struct {
	log *slog.Logger

	serviceTenderProvider ServiceTenderProvider
	serviceTenderCreator  ServiceTenderCreator
	serviceTenderEditor   ServiceTenderEditor

	serviceBidProvider      ServiceBidProvider
	serviceBidCreator       ServiceBidCreator
	serviceBidEditor        ServiceBidEditor
	serviceBidDecisionMaker ServiceBidDecisionMaker
	serviceBidFeedbacker    ServiceBidFeedbacker
}

func New(
	log *slog.Logger,

	serviceTenderProvider ServiceTenderProvider,
	serviceTenderCreator ServiceTenderCreator,
	serviceTenderEditor ServiceTenderEditor,

	serviceBidProvider ServiceBidProvider,
	serviceBidCreator ServiceBidCreator,
	serviceBidEditor ServiceBidEditor,
	serviceBidDecisionMaker ServiceBidDecisionMaker,
	serviceBidFeedbacker ServiceBidFeedbacker,

) *Api {
	return &Api{
		log: log,

		serviceTenderProvider: serviceTenderProvider,
		serviceTenderCreator:  serviceTenderCreator,
		serviceTenderEditor:   serviceTenderEditor,

		serviceBidProvider:      serviceBidProvider,
		serviceBidCreator:       serviceBidCreator,
		serviceBidEditor:        serviceBidEditor,
		serviceBidDecisionMaker: serviceBidDecisionMaker,
		serviceBidFeedbacker:    serviceBidFeedbacker,
	}
}

func (a *Api) Ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "ok")
}

func (a *Api) Placeholder(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "i'm a placeholder, please dont get attached to me:) ")
}
