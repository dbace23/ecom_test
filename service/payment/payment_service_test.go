package payment_test

import (
	"context"
	"errors"
	"testing"

	"ecom/model"
	paymentsvc "ecom/service/payment"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fakePaymentRepo struct {
	createCalled bool
	createInput  *model.Payment
	createErr    error
}

func (f *fakePaymentRepo) Create(ctx context.Context, p *model.Payment) error {
	f.createCalled = true
	f.createInput = p
	return f.createErr
}

func newServiceWithRepo(repo paymentsvc.Repository) paymentsvc.Service {
	return paymentsvc.NewService(repo)
}

func TestCreatePayment_SuccessAmountPositive(t *testing.T) {
	repo := &fakePaymentRepo{}
	svc := newServiceWithRepo(repo)

	req := model.CreatePaymentRequest{
		TransactionID: primitive.NewObjectID().Hex(),
		Amount:        100_000,
		Email:         "user@example.com",
	}

	p, err := svc.CreatePayment(context.Background(), req)
	if err != nil {
		t.Fatalf("CreatePayment returned error: %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repo.Create to be called")
	}

	if p.Status != model.PaymentStatusSuccess {
		t.Fatalf("expected payment status SUCCESS, got %s", p.Status)
	}

	if repo.createInput.Amount != req.Amount {
		t.Fatalf("expected repo payment amount %f, got %f", req.Amount, repo.createInput.Amount)
	}
}

func TestCreatePayment_FailedWhenAmountNonPositive(t *testing.T) {
	repo := &fakePaymentRepo{}
	svc := newServiceWithRepo(repo)

	req := model.CreatePaymentRequest{
		TransactionID: primitive.NewObjectID().Hex(),
		Amount:        0,
		Email:         "user@example.com",
	}

	p, err := svc.CreatePayment(context.Background(), req)
	if err != nil {
		t.Fatalf("CreatePayment returned error: %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repo.Create to be called")
	}

	if p.Status != model.PaymentStatusFailed {
		t.Fatalf("expected payment status FAILED, got %s", p.Status)
	}
}

func TestCreatePayment_RepoErrorIsReturned(t *testing.T) {
	repo := &fakePaymentRepo{
		createErr: errors.New("db error"),
	}
	svc := newServiceWithRepo(repo)

	req := model.CreatePaymentRequest{
		TransactionID: primitive.NewObjectID().Hex(),
		Amount:        50_000,
		Email:         "user@example.com",
	}

	_, err := svc.CreatePayment(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
