package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pozedorum/user-balance-service/internal/service"
)

type operationRoutes struct {
	operationService service.Operation
}

func newOperationRoutes(g *echo.Group, operationService service.Operation) {
	r := &operationRoutes{
		operationService: operationService,
	}

	// g.GET("/history", r.)
}

type getHistoryInput struct {
	AccountId int    `json:"account_id" validate:"required"`
	SortType  string `json:"sort_type,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

func (r *operationRoutes) getHistory(c echo.Context) error {
	var input getHistoryInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	operations, err := r.operationService.OperationHistory(c.Request().Context(), service.OperationHistoryInput{
		AccountId: input.AccountId,
		SortType:  input.SortType,
		Offset:    input.Offset,
		Limit:     input.Limit,
	})
	if err != nil {
		log.Debugf("operationRoutes.getHistory - r.operationService.OperationHistory: %v", err)
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
	}

	type responce struct {
		Operations []service.OperationHistoryOutput `json:"operations"`
	}

	return c.JSON(http.StatusOK, responce{
		Operations: operations,
	}) // TODO: add pagination
}

type getReportInput struct {
	Month int `json:"month" validate:"required"`
	Year  int `json:"year" validate:"required"`
}

func (r *operationRoutes) getReportLink(c echo.Context) error {
	var input getReportInput

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}
}
