package main

import (
	"log"

	controller "ecom/app/echoServer/controller"
	"ecom/app/echoServer/router"
	"ecom/config"
	paymentrepo "ecom/repository/payment"
	paymentservice "ecom/service/payment"
	"ecom/util/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load config (env)
	cfg := config.Load()

	// Connect ke Mongo
	client := database.NewMongoClient(cfg)
	paymentCol := database.PaymentCollection(client, cfg)

	// Wiring: repo → service → controller
	paymentRepo := paymentrepo.NewRepository(paymentCol)
	paymentSvc := paymentservice.NewService(paymentRepo)
	paymentCtrl := controller.NewPaymentController(paymentSvc)

	// Setup Echo
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//routes khusus Payment
	router.RegisterPaymentRoutes(e, paymentCtrl)

	log.Printf("Payment service listening on %s", cfg.PaymentPort)
	if err := e.Start(cfg.PaymentPort); err != nil {
		log.Fatal(err)
	}
}
