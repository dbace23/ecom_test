package controller

import (
	"net/http"

	"ecom/model"
	paymentservice "ecom/service/payment"

	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	svc paymentservice.Service
}

func NewPaymentController(svc paymentservice.Service) *PaymentController {
	return &PaymentController{svc: svc}
}

func (h *PaymentController) CreatePayment(c echo.Context) error {
	var req model.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid request body", err.Error())
	}

	if req.Amount <= 0 {
		return respondError(c, http.StatusBadRequest, "amount must be > 0", nil)
	}

	payment, err := h.svc.CreatePayment(c.Request().Context(), req)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, "failed to create payment", err.Error())
	}

	return respondOK(c, payment)
}
