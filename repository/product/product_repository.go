package product

import (
	"context"
	"time"

	"ecom/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, p *model.Product) error
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Product, error)
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type mongoRepository struct {
	col *mongo.Collection
}

func NewRepository(col *mongo.Collection) Repository {
	return &mongoRepository{col: col}
}

func (r *mongoRepository) Create(ctx context.Context, p *model.Product) error {
	p.ID = primitive.NewObjectID()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err := r.col.InsertOne(ctx, p)
	return err
}

func (r *mongoRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []model.Product
	if err := cur.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *mongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Product, error) {
	var p model.Product
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *mongoRepository) Update(ctx context.Context, p *model.Product) error {
	p.UpdatedAt = time.Now()
	_, err := r.col.UpdateByID(ctx, p.ID, bson.M{
		"$set": bson.M{
			"name":       p.Name,
			"price":      p.Price,
			"stock":      p.Stock,
			"updated_at": p.UpdatedAt,
		},
	})
	return err
}

func (r *mongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
