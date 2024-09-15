package api

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"zadanie-6105/internal/domain/model"
	"zadanie-6105/internal/domain/request"
	"zadanie-6105/internal/lib/binder"
	sl "zadanie-6105/internal/lib/logger/slog"
	"zadanie-6105/internal/lib/validator"
)

type ServiceBidProvider interface {
	GetBidsByUser(
		limit int32,
		offset int32,
		username string,
	) ([]model.BidResponse, error)
	BidsForTender(
		tenderId string,
		limit int32,
		offset int32,
		username string,
	) ([]model.BidResponse, error)
	BidStatus(
		bidId string,
		username string,
	) (string, error)
}
type ServiceBidCreator interface {
	CreateBid(
		name string,
		description string,
		tenderId string,
		authorType string,
		authorId string,
	) (model.BidResponse, error)
}
type ServiceBidEditor interface {
	UpdateBidStatus(
		bidId string,
		status string,
		username string,
	) (model.BidResponse, error)
	EditBid(
		bidId string,
		username string,
		name string,
		description string,
	) (model.BidResponse, error)
	RollbackBid(
		bidId string,
		version int32,
		username string,
	) (model.BidResponse, error)
}
type ServiceBidDecisionMaker interface {
	SubmitDecision(
		bidId string,
		decision string,
		username string,
	) (model.BidResponse, error)
}
type ServiceBidFeedbacker interface {
	Feedback(
		bidId string,
		feedback string,
		username string,
	) (model.BidResponse, error)
	Reviews(
		tenderId string,
		authorUsername string,
		requesterUsername string,
		limit int32,
		offset int32,
	) ([]model.Feedback, error)
}

func (a *Api) CreateBid(ctx echo.Context) error {
	const op = "Api.CreateBid"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.CreateBid{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bid model.BidResponse
	bid, err = a.serviceBidCreator.CreateBid(req.Name, req.Description, req.TenderId, req.AuthorType, req.AuthorId)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bid)
}

func (a *Api) GetBidsByUser(ctx echo.Context) error {
	const op = "Api.GetBidByUser"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.GetBidsByUser{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bids []model.BidResponse
	bids, err = a.serviceBidProvider.GetBidsByUser(req.Limit, req.Offset, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bids)
}

func (a *Api) BidsForTender(ctx echo.Context) error {
	const op = "Api.BidsForTender"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.BidsForTender{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bids []model.BidResponse
	bids, err = a.serviceBidProvider.BidsForTender(req.TenderId, req.Limit, req.Offset, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bids)
}

func (a *Api) BidStatus(ctx echo.Context) error {
	const op = "Api.BidStatus"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.BidStatus{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var status string
	status, err = a.serviceBidProvider.BidStatus(req.BidId, req.Username)
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, status)
}

func (a *Api) UpdateBidStatus(ctx echo.Context) error {
	const op = "Api.UpdateBidStatus"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.UpdateBidStatus{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bid model.BidResponse
	bid, err = a.serviceBidEditor.UpdateBidStatus(req.BidId, req.Status, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bid)
}

func (a *Api) EditBid(ctx echo.Context) error {
	const op = "Api.EditBid"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.EditBid{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bid model.BidResponse
	bid, err = a.serviceBidEditor.EditBid(req.BidId, req.Username, req.Name, req.Description)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bid)
}

func (a *Api) SubmitDecision(ctx echo.Context) error {
	const op = "Api.SubmitDecision"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.SubmitDecision{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bid model.BidResponse
	bid, err = a.serviceBidDecisionMaker.SubmitDecision(req.BidId, req.Decision, req.Username)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bid)
}

func (a *Api) Feedback(ctx echo.Context) error {
	const op = "Api.Feedback"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.Feedback{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bid model.BidResponse
	bid, err = a.serviceBidFeedbacker.Feedback(req.BidId, req.BidFeedback, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bid)
}

func (a *Api) RollbackBid(ctx echo.Context) error {
	const op = "Api.RollbackBid"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.RollbackBid{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var bid model.BidResponse
	bid, err = a.serviceBidEditor.RollbackBid(req.BidId, req.Version, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, bid)
}

func (a *Api) Reviews(ctx echo.Context) error {
	const op = "Api.Reviews"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.Reviews{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var feedbacks []model.Feedback
	feedbacks, err = a.serviceBidFeedbacker.Reviews(req.TenderId, req.AuthorUsername, req.RequesterUsername, req.Limit, req.Offset)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, feedbacks)
}
