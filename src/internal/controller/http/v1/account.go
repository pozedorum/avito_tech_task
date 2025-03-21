package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pozedorum/user-balance-service/internal/service"
)

type accountRoutes struct {
	accountService service.Account
}

func NewAccountRoutes (g *echo.Group,accountService service.Account) {
	r := &accountRoutes{
		accountService: accountService,
	}
	
	g.POST("/create",r.create)
	g.POST("/deposit", r.deposit)
	g.POST("/withdraw", r.withdraw)
	g.POST("/transfer", r.transfer)
	g.POST("/", r.getBalance)
}


func (r *accountRoutes) create (c echo.Context) error {
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

	return c.JSON(http.StatusCreated, responce {Id: id})
}

func (r *accountRoutes) deposit (c echo.Context) error {
	var input service.AccountDepositInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request error")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Deposit(c.Request().Context(), service.AccountDepositInput{
		Id: input.Id,
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

	return c.JSON(http.StatusOK, map[string]interface{}{"message":"success"})
}

func (r *accountRoutes) withdraw (c echo.Context) error {
	var input service.AccountWithDrawInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request error")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Withdraw(c.Request().Context(), service.AccountWithDrawInput{
		Id: input.Id,
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

	return c.JSON(http.StatusOK, map[string]interface{}{"message":"success"})
}

func (r *accountRoutes) transfer (c echo.Context) error {
	var input service.AccountTransferInput

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
		ToId: input.ToId,
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

	return c.JSON(http.StatusOK, map[string]interface{}{"message":"success"})
}

func (r *accountRoutes) getBalance (c echo.Context) error {
	var input service.

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request error")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.GetAccountById(c.Request().Context(), service.AccountTransferInput{
		FromId: input.FromId,
		ToId: input.ToId,
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

	return c.JSON(http.StatusOK, map[string]interface{}{"message":"success"})
}

// Поменять все input структуры на новые, новые написать с учётом тегов