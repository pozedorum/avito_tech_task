package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pozedorum/user-balance-service/internal/service"
)

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{
		authService: authService,
	}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
}

type signUpInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

func (r *authRoutes) signUp(c echo.Context) error {
	var input signUpInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == service.ErrUserAlreadyExists {
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

type signInInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

func (r *authRoutes) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == service.ErrUserNotFound {
			newErrResponce(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusCreated, response{Token: token})
}
