package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ecom/model"
	txsvc "ecom/service/transaction"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fakeProductRepo struct {
	findByIDResult *model.Product
	findByIDErr    error

	updateCalled bool
	updateInput  *model.Product
	updateErr    error
}

func (f *fakeProductRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Product, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeProductRepo) Update(ctx context.Context, p *model.Product) error {
	f.updateCalled = true
	f.updateInput = p
	return f.updateErr
}

type fakeTxRepo struct {
	createCalled bool
	createInput  *model.Transaction
	createErr    error

	findAllResult []model.Transaction
	findAllErr    error

	findByIDResult *model.Transaction
	findByIDErr    error

	updateCalled bool
	updateInput  *model.Transaction
	updateErr    error

	deleteCalled bool
	deleteID     primitive.ObjectID
	deleteErr    error

	expireCalled    bool
	expireOlderThan time.Duration
	expireResult    int64
	expireErr       error
}

func (f *fakeTxRepo) Create(ctx context.Context, t *model.Transaction) error {
	f.createCalled = true
	f.createInput = t

	if t.ID.IsZero() {
		t.ID = primitive.NewObjectID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return f.createErr
}

func (f *fakeTxRepo) FindAll(ctx context.Context) ([]model.Transaction, error) {
	return f.findAllResult, f.findAllErr
}

func (f *fakeTxRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Transaction, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeTxRepo) Update(ctx context.Context, t *model.Transaction) error {
	f.updateCalled = true
	f.updateInput = t
	return f.updateErr
}

func (f *fakeTxRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	f.deleteCalled = true
	f.deleteID = id
	return f.deleteErr
}

func (f *fakeTxRepo) ExpireOldPending(ctx context.Context, olderThan time.Duration) (int64, error) {
	f.expireCalled = true
	f.expireOlderThan = olderThan
	return f.expireResult, f.expireErr
}

type fakePaymentClient struct {
	resp *model.Payment
	err  error

	called bool
	input  model.CreatePaymentRequest
}

func (f *fakePaymentClient) CreatePayment(ctx context.Context, req model.CreatePaymentRequest) (*model.Payment, error) {
	f.called = true
	f.input = req
	return f.resp, f.err
}

func newService(
	prodRepo txsvc.ProductRepository,
	txRepo txsvc.TransactionRepository,
	payment txsvc.PaymentClient,
) txsvc.Service {
	return txsvc.NewService(prodRepo, txRepo, payment)
}

func TestCreateTransaction_SuccessPaymentSuccess(t *testing.T) {

	productID := primitive.NewObjectID()
	product := &model.Product{
		ID:    productID,
		Name:  "Lapangan Futsal",
		Price: 100_000,
		Stock: 10,
	}

	prodRepo := &fakeProductRepo{
		findByIDResult: product,
	}

	txRepo := &fakeTxRepo{}
	paymentClient := &fakePaymentClient{
		resp: &model.Payment{
			Status: model.PaymentStatusSuccess,
		},
	}

	svc := newService(prodRepo, txRepo, paymentClient)

	req := model.CreateTransactionRequest{
		ProductID: productID.Hex(),
		Qty:       2,
		Email:     "user@example.com",
	}

	tx, err := svc.CreateTransaction(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateTransaction returned error: %v", err)
	}

	if !txRepo.createCalled {
		t.Fatal("expected txRepo.Create to be called")
	}
	if !paymentClient.called {
		t.Fatal("expected paymentClient.CreatePayment to be called")
	}
	if !txRepo.updateCalled {
		t.Fatal("expected txRepo.Update to be called")
	}
	if !prodRepo.updateCalled {
		t.Fatal("expected productRepo.Update to be called when payment success")
	}

	if tx.Status != model.TransactionStatusSuccess {
		t.Fatalf("expected transaction status SUCCESS, got %s", tx.Status)
	}
	if prodRepo.updateInput.Stock != 8 { // 10 - 2
		t.Fatalf("expected product stock 8, got %d", prodRepo.updateInput.Stock)
	}
}

func TestCreateTransaction_InsufficientStock(t *testing.T) {
	productID := primitive.NewObjectID()
	product := &model.Product{
		ID:    productID,
		Name:  "Lapangan Futsal",
		Price: 100_000,
		Stock: 1, // stok cuma 1
	}

	prodRepo := &fakeProductRepo{
		findByIDResult: product,
	}

	txRepo := &fakeTxRepo{}
	paymentClient := &fakePaymentClient{}

	svc := newService(prodRepo, txRepo, paymentClient)

	req := model.CreateTransactionRequest{
		ProductID: productID.Hex(),
		Qty:       2, // butuh 2
		Email:     "user@example.com",
	}

	_, err := svc.CreateTransaction(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for insufficient stock, got nil")
	}

	if txRepo.createCalled {
		t.Fatal("expected txRepo.Create NOT to be called when stock insufficient")
	}
	if paymentClient.called {
		t.Fatal("expected paymentClient NOT to be called when stock insufficient")
	}
}

func TestCreateTransaction_PaymentErrorMarksFailed(t *testing.T) {
	productID := primitive.NewObjectID()
	product := &model.Product{
		ID:    productID,
		Name:  "Lapangan Futsal",
		Price: 100_000,
		Stock: 10,
	}

	prodRepo := &fakeProductRepo{
		findByIDResult: product,
	}

	txRepo := &fakeTxRepo{}
	paymentClient := &fakePaymentClient{
		err: errors.New("payment service down"),
	}

	svc := newService(prodRepo, txRepo, paymentClient)

	req := model.CreateTransactionRequest{
		ProductID: productID.Hex(),
		Qty:       2,
		Email:     "user@example.com",
	}

	_, err := svc.CreateTransaction(context.Background(), req)
	if err == nil {
		t.Fatal("expected error when payment client fails, got nil")
	}

	if !txRepo.updateCalled {
		t.Fatal("expected txRepo.Update called to mark FAILED")
	}
	if txRepo.updateInput.Status != model.TransactionStatusFailed {
		t.Fatalf("expected transaction status FAILED, got %s", txRepo.updateInput.Status)
	}
}

func TestRunExpireJob_CallsRepoWithDuration(t *testing.T) {
	prodRepo := &fakeProductRepo{}
	txRepo := &fakeTxRepo{
		expireResult: 5,
	}
	paymentClient := &fakePaymentClient{}

	svc := newService(prodRepo, txRepo, paymentClient)

	ctx := context.Background()
	modified, err := svc.RunExpireJob(ctx)
	if err != nil {
		t.Fatalf("RunExpireJob returned error: %v", err)
	}

	if !txRepo.expireCalled {
		t.Fatal("expected ExpireOldPending to be called")
	}
	if modified != 5 {
		t.Fatalf("expected modified = 5, got %d", modified)
	}
}
