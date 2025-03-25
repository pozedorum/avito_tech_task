package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pozedorum/user-balance-service/internal/service"
)

type accountRoutes struct {
	accountService service.Account
}

func NewAccountRoutes(g *echo.Group, accountService service.Account) {
	r := &accountRoutes{
		accountService: accountService,
	}

	g.POST("/create", r.create)
	g.POST("/deposit", r.deposit)
	g.POST("/withdraw", r.withdraw)
	g.POST("/transfer", r.transfer)
	g.POST("/", r.getBalance)

}

func (r *accountRoutes) create(c echo.Context) error {
	id, err := r.accountService.CreateAccount(c.Request().Context())
	if err != nil {
		if err == service.ErrAccountAlreadyExists {
			newErrResponce(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type responce struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, responce{Id: id})
}

type accountDepositInput struct {
	Id     int `json:"id" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}

func (r *accountRoutes) deposit(c echo.Context) error {
	var input accountDepositInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request error")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Deposit(c.Request().Context(), service.AccountDepositInput{
		Id:     input.Id,
		Amount: input.Amount,
	})

	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrResponce(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success"})
}

type accountWithDrawInput struct {
	Id     int `json:"id" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}

func (r *accountRoutes) withdraw(c echo.Context) error {
	var input accountWithDrawInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request error")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Withdraw(c.Request().Context(), service.AccountWithDrawInput{
		Id:     input.Id,
		Amount: input.Amount,
	})

	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrResponce(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrResponce(c, http.StatusBadRequest, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success"})
}

type accountTransferInput struct {
	FromId int `json:"from" validate:"required"`
	ToId   int `json:"to" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}

func (r *accountRoutes) transfer(c echo.Context) error {
	var input accountTransferInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request error")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Transfer(c.Request().Context(), service.AccountTransferInput{
		FromId: input.FromId,
		ToId:   input.ToId,
		Amount: input.Amount,
	})

	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrResponce(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrResponce(c, http.StatusBadRequest, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success"})
}

type getBalanceInput struct {
	Id int `json:"id" validate:"required"`
}

func (r *accountRoutes) getBalance(c echo.Context) error {
	var input getBalanceInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request error")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	account, err := r.accountService.GetAccountById(c.Request().Context(), input.Id)
	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrResponce(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrResponce(c, http.StatusBadRequest, "internal server error")
		return err
	}

	type responce struct {
		Id      int `json:"id"`
		Balance int `json:"balance"`
	}

	return c.JSON(http.StatusOK, responce{
		Id:      account.Id,
		Balance: account.Balance,
	})
}
