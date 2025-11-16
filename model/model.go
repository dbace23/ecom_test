package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string

const (
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionID primitive.ObjectID `bson:"transaction_id" json:"transaction_id"`
	Amount        float64            `bson:"amount" json:"amount"`
	Email         string             `bson:"email" json:"email"`
	Status        PaymentStatus      `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type CreatePaymentRequest struct {
	TransactionID string  `json:"transaction_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Email         string  `json:"email" validate:"required,email"`
}

type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Price     float64            `bson:"price" json:"price"`
	Stock     int                `bson:"stock" json:"stock"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required,gt=0"`
	Stock int     `json:"stock" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required,gt=0"`
	Stock int     `json:"stock" validate:"required,gte=0"`
}

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusFailed  TransactionStatus = "FAILED"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	Qty         int                `bson:"qty" json:"qty"`
	TotalAmount float64            `bson:"total_amount" json:"total_amount"`
	Email       string             `bson:"email" json:"email"`
	Status      TransactionStatus  `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateTransactionRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required,gt=0"`
	Email     string `json:"email" validate:"required,email"`
}

type UpdateTransactionRequest struct {
	Qty   int    `json:"qty" validate:"required,gt=0"`
	Email string `json:"email" validate:"required,email"`
}
