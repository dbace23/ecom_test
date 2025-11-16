package payment

import (
	"context"
	"time"

	"ecom/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, p *model.Payment) error
}

type repo struct {
	col *mongo.Collection
}

func NewRepository(col *mongo.Collection) Repository {
	return &repo{col: col}
}

func (r *repo) Create(ctx context.Context, p *model.Payment) error {
	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, p)
	return err
}
