package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pozedorum/user-balance-service/internal/entity"
	"github.com/pozedorum/user-balance-service/internal/service"
)

type productRoutes struct {
	productService service.Product
}

func newProductRoutes(g *echo.Group, productService service.Product) {
	r := &productRoutes{
		productService: productService,
	}

	g.POST("/create", r.create)
	g.GET("/", r.getById)
}

type productCreateInput struct {
	Name string `json:"name" validate:"required"`
}

func (r *productRoutes) create(c echo.Context) error {
	var input productCreateInput

	id, err := r.productService.CreateProduct(c.Request().Context(), input.Name)
	if err != nil {
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type responce struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, responce{Id: id})
}

type getByIdInput struct {
	Id int `json:"id" validate:"required"`
}

func (r *productRoutes) getById(c echo.Context) error {
	var input getByIdInput

	if err := c.Bind(&input); err != nil {
		newErrResponce(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrResponce(c, http.StatusBadRequest, err.Error())
		return err
	}

	product, err := r.productService.GetProductById(c.Request().Context(), input.Id)
	if err != nil {
		newErrResponce(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Product entity.Product `json:"product"`
	}

	return c.JSON(http.StatusOK, response{
		Product: product,
	})
}
