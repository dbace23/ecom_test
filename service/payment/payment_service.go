package payment

import (
	"context"
	"fmt"

	"ecom/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CreatePayment(ctx context.Context, req model.CreatePaymentRequest) (*model.Payment, error)
}

type Repository interface {
	Create(ctx context.Context, p *model.Payment) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePayment(ctx context.Context, req model.CreatePaymentRequest) (*model.Payment, error) {
	txID, err := primitive.ObjectIDFromHex(req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction_id: %w", err)
	}

	status := model.PaymentStatusSuccess
	if req.Amount <= 0 {
		status = model.PaymentStatusFailed
	}

	p := &model.Payment{
		TransactionID: txID,
		Amount:        req.Amount,
		Email:         req.Email,
		Status:        status,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
