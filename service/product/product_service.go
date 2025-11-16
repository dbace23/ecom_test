package product

import (
	"context"
	"fmt"

	"ecom/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(ctx context.Context, p *model.Product) error
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Product, error)
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type Service interface {
	Create(ctx context.Context, req model.CreateProductRequest) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id string) (*model.Product, error)
	Update(ctx context.Context, id string, req model.UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req model.CreateProductRequest) (*model.Product, error) {
	p := &model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (*model.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}
	return s.repo.FindByID(ctx, objID)
}

func (s *service) Update(ctx context.Context, id string, req model.UpdateProductRequest) (*model.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}
	p, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	p.Name = req.Name
	p.Price = req.Price
	p.Stock = req.Stock

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	return s.repo.Delete(ctx, objID)
}
