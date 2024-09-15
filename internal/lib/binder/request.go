package binder

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	sl "zadanie-6105/internal/lib/logger/slog"
)

func BindReq(log *slog.Logger, ctx echo.Context, req interface{}) error {
	err := ctx.Bind(req)
	if err != nil {
		log.Error("failed to parse", sl.Err(err))
		return err
	}

	binder := &echo.DefaultBinder{}
	err = binder.BindHeaders(ctx, req)
	if err != nil {
		log.Error("failed to parse token", sl.Err(err))
		return err
	}

	err = binder.BindPathParams(ctx, req)
	if err != nil {
		log.Error("failed to parse path", sl.Err(err))
		return err
	}
	err = binder.BindQueryParams(ctx, req)
	if err != nil {
		log.Error("failed to parse query", sl.Err(err))
		return err
	}

	return nil
}
