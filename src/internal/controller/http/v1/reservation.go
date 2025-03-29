package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pozedorum/user-balance-service/internal/service"
)

type reservationRoutes struct {
	reservationService service.Reservation
}

func newReservationRoutes(g *echo.Group, reservationService service.Reservation) {
	r := reservationRoutes{
		reservationService: reservationService,
	}

	g.POST("/create", r.create)
	g.POST("/revenue", r.revenue)
	g.POST("/refund", r.refund)
}

type reservationCreateInput struct {
	AccountId int `json:"account_id" validate:"required"`
	ProductId int `json:"product_id" validate:"required"`
	OrderId   int `json:"order_id" validate:"required"`
	Amount    int `json:"amount" validate:"required"`
}

func (r *reservationRoutes) create(c echo.Context) error {
	var input reservationCreateInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.reservationService.CreateReservation(c.Request().Context(), service.ReservationCreateInput{
		AccountId: input.AccountId,
		ProductId: input.ProductId,
		OrderId:   input.OrderId,
		Amount:    input.Amount,
	})
	if err != nil {
		if err == service.ErrCannotCreateReservation {
			newErrResponce(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{Id: id})
}

type reservationRevenueInput struct {
	AccountId int `json:"account_id" validate:"required"`
	ProductId int `json:"product_id" validate:"required"`
	OrderId   int `json:"order_id" validate:"required"`
	Amount    int `json:"amount" validate:"required"`
}

func (r *reservationRoutes) revenue(c echo.Context) error {
	var input reservationRevenueInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.reservationService.RevenueReservationById(c.Request().Context(), input.AccountId)
	if err != nil {
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success"})
}

type reservationRefundInput struct {
	OrderId int `json:"order_id" validate:"required"`
}

func (r *reservationRoutes) refund(c echo.Context) error {
	var input reservationRefundInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"mesage": "success"})
}
