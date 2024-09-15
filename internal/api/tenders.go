package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
	"zadanie-6105/internal/domain/model"
	"zadanie-6105/internal/domain/request"
	"zadanie-6105/internal/lib/binder"
	sl "zadanie-6105/internal/lib/logger/slog"
	"zadanie-6105/internal/lib/validator"
)

type ServiceTenderProvider interface {
	Tenders(
		limit int32,
		offset int32,
		serviceType []string,
	) ([]model.TenderResponse, error)
	GetTenderByUser(
		limit int32,
		offset int32,
		username string,
	) ([]model.TenderResponse, error)
	TenderStatus(
		tenderId string,
		username string,
	) (string, error)
}
type ServiceTenderCreator interface {
	CreateTender(
		name string,
		description string,
		serviceType string,
		organizationId string,
		creatorUsername string,
	) (model.TenderResponse, error)
}
type ServiceTenderEditor interface {
	ChangeTenderStatus(
		tenderId string,
		status string,
		username string,
	) (model.TenderResponse, error)
	EditTender(
		tenderId string,
		username string,
		name string,
		description string,
		serviceType string,
	) (model.TenderResponse, error)
	RollbackTender(
		tenderId string,
		version int32,
		username string,
	) (model.TenderResponse, error)
}

func (a *Api) Tenders(ctx echo.Context) error {
	const op = "Api.Tenders"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.GetTender{}

	err := binder.BindReq(log, ctx, &req)
	if len(req.ServiceType) != 0 {
		req.ServiceType = strings.Split(req.ServiceType[0], ",")
	}
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	fmt.Println(req.ServiceType)
	var tenders []model.TenderResponse
	tenders, err = a.serviceTenderProvider.Tenders(req.Limit, req.Offset, req.ServiceType)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, tenders)
}

func (a *Api) CreateTender(ctx echo.Context) error {
	const op = "Api.CreateTender"

	log := a.log.With(
		slog.String("op", op),
	)
	req := request.CreateTender{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var tender model.TenderResponse
	tender, err = a.serviceTenderCreator.CreateTender(req.Name, req.Description, req.ServiceType, req.OrganizationId, req.CreatorUsername)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, tender)
}

func (a *Api) GetTenderByUser(ctx echo.Context) error {
	const op = "Api.GetTenderByUser"

	log := a.log.With(
		slog.String("op", op),
	)
	req := request.GetTenderByUser{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var tenders []model.TenderResponse
	tenders, err = a.serviceTenderProvider.GetTenderByUser(req.Limit, req.Offset, req.Username)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, tenders)
}

func (a *Api) TenderStatus(ctx echo.Context) error {
	const op = "Api.TenderStatus"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.TenderStatus{}
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
	status, err = a.serviceTenderProvider.TenderStatus(req.TenderId, req.Username)

	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, status)
}

func (a *Api) ChangeTenderStatus(ctx echo.Context) error {
	const op = "Api.ChangeTenderStatus"
	log := a.log.With(
		slog.String("op", op),
	)

	req := request.UpdateTenderStatus{}
	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var tender model.TenderResponse
	tender, err = a.serviceTenderEditor.ChangeTenderStatus(req.TenderId, req.Status, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, tender)
}

func (a *Api) EditTender(ctx echo.Context) error {
	const op = "Api.EditTender"

	log := a.log.With(
		slog.String("op", op),
	)
	req := request.EditTender{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	err = validator.Validate(req)
	if err != nil {
		return err
	}

	var tender model.TenderResponse
	tender, err = a.serviceTenderEditor.EditTender(req.TenderId, req.Username, req.Name, req.Description, req.ServiceType)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, tender)
}

func (a *Api) RollbackTender(ctx echo.Context) error {
	const op = "Api.RollbackTender"

	log := a.log.With(
		slog.String("op", op),
	)
	req := request.RollbackTender{}

	err := binder.BindReq(log, ctx, &req)
	if err != nil {
		return err
	}
	err = validator.Validate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Info(sl.Req(req))

	var tender model.TenderResponse
	tender, err = a.serviceTenderEditor.RollbackTender(req.TenderId, req.Version, req.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, tender)
}
