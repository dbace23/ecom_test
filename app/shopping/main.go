package main

import (
	"log"

	"ecom/app/cron/shopping"
	"ecom/app/echoServer/controller"
	"ecom/app/echoServer/router"
	"ecom/config"
	productrepo "ecom/repository/product"
	txrepo "ecom/repository/transaction"
	productservice "ecom/service/product"
	txservice "ecom/service/transaction"
	"ecom/util/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect Mongo
	client := database.NewMongoClient(cfg)
	productCol := database.ProductCollection(client, cfg)
	txCol := database.TransactionCollection(client, cfg)

	//Repo
	prodRepo := productrepo.NewRepository(productCol)
	transactionRepo := txrepo.NewRepository(txCol)

	// Payment client
	paymentClient := txservice.NewHTTPPaymentClient(cfg.PaymentBaseURL)

	// Service
	prodSvc := productservice.NewService(prodRepo)
	txSvc := txservice.NewService(prodRepo, transactionRepo, paymentClient)

	//Start cron job
	shopping.StartTransactionExpireJob(txSvc)

	// Echo & controllers
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	productCtrl := controller.NewProductController(prodSvc)
	transactionCtrl := controller.NewTransactionController(txSvc)

	//routes shopping (products + transactions)
	router.RegisterShoppingRoutes(e, productCtrl, transactionCtrl)

	log.Printf("Shopping service listening on %s", cfg.ShoppingPort)
	if err := e.Start(cfg.ShoppingPort); err != nil {
		log.Fatal(err)
	}
}
