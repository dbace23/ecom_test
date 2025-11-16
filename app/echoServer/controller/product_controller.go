package controller

import (
	"net/http"

	"ecom/model"
	productservice "ecom/service/product"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	svc productservice.Service
}

func NewProductController(svc productservice.Service) *ProductController {
	return &ProductController{svc: svc}
}

func (h *ProductController) Create(c echo.Context) error {
	var req model.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body", err.Error())
	}

	p, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to create product", err.Error())
	}

	return respondOK(c, p)
}

func (h *ProductController) GetAll(c echo.Context) error {
	products, err := h.svc.GetAll(c.Request().Context())
	if err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to get products", err.Error())
	}
	return respondOK(c, products)
}

func (h *ProductController) GetByID(c echo.Context) error {
	id := c.Param("id")
	p, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		return respondError(c, http.StatusNotFound, "product not found", err.Error())
	}
	return respondOK(c, p)
}

func (h *ProductController) Update(c echo.Context) error {
	id := c.Param("id")
	var req model.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body", err.Error())
	}

	p, err := h.svc.Update(c.Request().Context(), id, req)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to update product", err.Error())
	}
	return respondOK(c, p)
}

func (h *ProductController) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to delete product", err.Error())
	}
	return respondOK(c, echo.Map{"deleted": true})
}
