package router

import (
	Controller "ecom/app/echoServer/controller"

	"github.com/labstack/echo/v4"
)

func New(e *echo.Echo, paymentController *Controller.PaymentController) {
	e.POST("/payments", paymentController.CreatePayment)
}

func RegisterPaymentRoutes(e *echo.Echo, paymentController *Controller.PaymentController) {
	e.POST("/payments", paymentController.CreatePayment)
}

func RegisterShoppingRoutes(
	e *echo.Echo,
	productController *Controller.ProductController,
	transactionController *Controller.TransactionController,
) {
	// products
	e.POST("/products", productController.Create)
	e.GET("/products", productController.GetAll)
	e.GET("/products/:id", productController.GetByID)
	e.PUT("/products/:id", productController.Update)
	e.DELETE("/products/:id", productController.Delete)

	// transactions
	e.POST("/transactions", transactionController.Create)
	e.GET("/transactions", transactionController.GetAll)
	e.GET("/transactions/:id", transactionController.GetByID)
	e.PUT("/transactions/:id", transactionController.Update)
	e.DELETE("/transactions/:id", transactionController.Delete)
}
