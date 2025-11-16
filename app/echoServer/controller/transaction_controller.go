package controller

import (
	"net/http"

	"ecom/model"
	txservice "ecom/service/transaction"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	svc txservice.Service
}

func NewTransactionController(svc txservice.Service) *TransactionController {
	return &TransactionController{svc: svc}
}

func (h *TransactionController) Create(c echo.Context) error {
	var req model.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body", err.Error())
	}

	tx, err := h.svc.CreateTransaction(c.Request().Context(), req)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to create transaction", err.Error())
	}

	return respondOK(c, tx)
}

func (h *TransactionController) GetAll(c echo.Context) error {
	txs, err := h.svc.GetAll(c.Request().Context())
	if err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to get transactions", err.Error())
	}
	return respondOK(c, txs)
}

func (h *TransactionController) GetByID(c echo.Context) error {
	id := c.Param("id")
	tx, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		return respondError(c, http.StatusNotFound, "transaction not found", err.Error())
	}
	return respondOK(c, tx)
}

func (h *TransactionController) Update(c echo.Context) error {
	id := c.Param("id")
	var req model.UpdateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body", err.Error())
	}

	tx, err := h.svc.Update(c.Request().Context(), id, req)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to update transaction", err.Error())
	}
	return respondOK(c, tx)
}

func (h *TransactionController) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to delete transaction", err.Error())
	}
	return respondOK(c, echo.Map{"deleted": true})
}
