package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ecom/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CreateTransaction(ctx context.Context, req model.CreateTransactionRequest) (*model.Transaction, error)
	GetAll(ctx context.Context) ([]model.Transaction, error)
	GetByID(ctx context.Context, id string) (*model.Transaction, error)
	Update(ctx context.Context, id string, req model.UpdateTransactionRequest) (*model.Transaction, error)
	Delete(ctx context.Context, id string) error
	RunExpireJob(ctx context.Context) (int64, error)
}

type ProductRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Product, error)
	Update(ctx context.Context, p *model.Product) error
}

type TransactionRepository interface {
	Create(ctx context.Context, t *model.Transaction) error
	FindAll(ctx context.Context) ([]model.Transaction, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Transaction, error)
	Update(ctx context.Context, t *model.Transaction) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	ExpireOldPending(ctx context.Context, olderThan time.Duration) (int64, error)
}

type PaymentClient interface {
	CreatePayment(ctx context.Context, req model.CreatePaymentRequest) (*model.Payment, error)
}

type httpPaymentClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPPaymentClient(baseURL string) PaymentClient {
	return &httpPaymentClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *httpPaymentClient) CreatePayment(ctx context.Context, req model.CreatePaymentRequest) (*model.Payment, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := c.baseURL + "/payments"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("payment service returned status %d", resp.StatusCode)
	}

	var payment model.Payment
	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

type service struct {
	productRepo ProductRepository
	txRepo      TransactionRepository
	payment     PaymentClient
}

func NewService(
	productRepo ProductRepository,
	txRepo TransactionRepository,
	payment PaymentClient,
) Service {
	return &service{
		productRepo: productRepo,
		txRepo:      txRepo,
		payment:     payment,
	}
}

// /transactions (POST)
func (s *service) CreateTransaction(ctx context.Context, req model.CreateTransactionRequest) (*model.Transaction, error) {
	// Ambil & validasi product
	prodID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("invalid product_id")
	}

	prod, err := s.productRepo.FindByID(ctx, prodID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if prod.Stock < req.Qty {
		return nil, fmt.Errorf("insufficient stock")
	}

	total := prod.Price * float64(req.Qty)

	//  Buat transaksi PENDING
	tx := &model.Transaction{
		ProductID:   prod.ID,
		Qty:         req.Qty,
		TotalAmount: total,
		Email:       req.Email,
		Status:      model.TransactionStatusPending,
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	// Call Payment service
	payReq := model.CreatePaymentRequest{
		TransactionID: tx.ID.Hex(),
		Amount:        total,
		Email:         req.Email,
	}

	payment, err := s.payment.CreatePayment(ctx, payReq)
	if err != nil {
		// kalau error call payment  FAILED
		tx.Status = model.TransactionStatusFailed
		_ = s.txRepo.Update(ctx, tx)
		return nil, fmt.Errorf("payment error: %w", err)
	}

	//  Update status & stok berdasarkan hasil payment
	if payment.Status == model.PaymentStatusSuccess {
		tx.Status = model.TransactionStatusSuccess
		prod.Stock -= req.Qty

		if err := s.productRepo.Update(ctx, prod); err != nil {
			return nil, fmt.Errorf("update product stock: %w", err)
		}
	} else {
		tx.Status = model.TransactionStatusFailed
	}

	if err := s.txRepo.Update(ctx, tx); err != nil {
		return nil, fmt.Errorf("update transaction: %w", err)
	}

	return tx, nil
}

// /transactions (GET)
func (s *service) GetAll(ctx context.Context) ([]model.Transaction, error) {
	return s.txRepo.FindAll(ctx)
}

// /transactions/{id} (GET)
func (s *service) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}
	return s.txRepo.FindByID(ctx, objID)
}

// /transactions/{id} (PUT)
func (s *service) Update(ctx context.Context, id string, req model.UpdateTransactionRequest) (*model.Transaction, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	tx, err := s.txRepo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	tx.Qty = req.Qty
	tx.Email = req.Email
	tx.UpdatedAt = time.Now()

	if err := s.txRepo.Update(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

// /transactions/{id} (DELETE)
func (s *service) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	return s.txRepo.Delete(ctx, objID)
}

// cron job transaksi PENDING yang terlalu lama
func (s *service) RunExpireJob(ctx context.Context) (int64, error) {
	// expire PENDING lebih tua dari 30 menit
	return s.txRepo.ExpireOldPending(ctx, 30*time.Minute)
}
