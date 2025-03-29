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

func newOperationRoutes(g *echo.Group, operationService service.Operation) *operationRoutes {
	r := &operationRoutes{
		operationService: operationService,
	}

	g.GET("/history", r.getHistory)
	g.GET("/report-link", r.getReportLink)
	g.GET("/report-file", r.getReportFile)

	return r
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

	if err := c.Bind(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	link, err := r.operationService.MakeReportLink(c.Request().Context(), input.Month, input.Year)
	if err != nil {
		log.Debugf("operationRoutes.getReportLink - r.operationService.MakeReportLink: %v", err)
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Link string `json:"link"`
	}

	return c.JSON(http.StatusOK, response{Link: link})
}

func (r *operationRoutes) getReportFile(c echo.Context) error {
	var input getReportInput

	if err := c.Bind(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	file, err := r.operationService.MakeReportFile(c.Request().Context(), input.Month, input.Year)
	if err != nil {
		log.Debugf("operationRoutes.getReportFile - r.operationService.MakeReportFile: %v", err)
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.Blob(http.StatusOK, "text/csv", file)
}
