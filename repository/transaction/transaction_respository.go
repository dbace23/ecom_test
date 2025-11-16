package transaction

import (
	"context"
	"time"

	"ecom/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, t *model.Transaction) error
	FindAll(ctx context.Context) ([]model.Transaction, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Transaction, error)
	Update(ctx context.Context, t *model.Transaction) error
	Delete(ctx context.Context, id primitive.ObjectID) error

	ExpireOldPending(ctx context.Context, olderThan time.Duration) (int64, error)
}

type mongoRepository struct {
	col *mongo.Collection
}

func NewRepository(col *mongo.Collection) Repository {
	return &mongoRepository{col: col}
}

func (r *mongoRepository) Create(ctx context.Context, t *model.Transaction) error {
	t.ID = primitive.NewObjectID()
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	_, err := r.col.InsertOne(ctx, t)
	return err
}

func (r *mongoRepository) FindAll(ctx context.Context) ([]model.Transaction, error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var txs []model.Transaction
	if err := cur.All(ctx, &txs); err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *mongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Transaction, error) {
	var t model.Transaction
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *mongoRepository) Update(ctx context.Context, t *model.Transaction) error {
	t.UpdatedAt = time.Now()
	_, err := r.col.UpdateByID(ctx, t.ID, bson.M{
		"$set": bson.M{
			"product_id":   t.ProductID,
			"qty":          t.Qty,
			"total_amount": t.TotalAmount,
			"email":        t.Email,
			"status":       t.Status,
			"updated_at":   t.UpdatedAt,
		},
	})
	return err
}

func (r *mongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoRepository) ExpireOldPending(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)

	res, err := r.col.UpdateMany(ctx,
		bson.M{
			"status":     model.TransactionStatusPending,
			"created_at": bson.M{"$lt": cutoff},
		},
		bson.M{
			"$set": bson.M{
				"status":     model.TransactionStatusFailed,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
