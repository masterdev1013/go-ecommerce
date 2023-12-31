package product

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/mongodb"
)

// Repository db repository
type Repository interface {
	// get product
	Get(ctx context.Context, id string) (prod product, err error)

	// create new product
	Create(ctx context.Context, req createRequest) (id string, err error)

	// update product
	Update(ctx context.Context, id string, req updateRequest) error

	// delete product
	Delete(ctx context.Context, id string) error
}

type repository struct {
	mdb mongodb.Provider
}

// NewRepository returns a new repostory
func NewRepository(mdb mongodb.Provider) Repository {
	return repository{mdb}
}

func (r repository) Get(ctx context.Context, id string) (prod product, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return prod, err
	}

	var prodM bson.M
	if err := r.mdb.FindOne("product", ctx, bson.M{"_id": objectID}).Decode(&prodM); err != nil {
		return prod, errors.Wrapf(errorsd.New("error getting product in db"), "error getting product in db: %s", err.Error())
	}
	bsonBytes, _ := bson.Marshal(prodM)
	bson.Unmarshal(bsonBytes, &prod)
	return prod, err
}

func (r repository) Create(ctx context.Context, req createRequest) (id string, err error) {
	// @todo add product/category validation
	res, err := r.mdb.InsertOne("product", ctx, bson.D{
		{"title", req.Title},
		{"body", req.Body},
		{"category", req.Category},
		{"assets", req.Assets},
		{"variants", req.Variants}})
	if err != nil {
		return id, errors.Wrapf(errorsd.New("error inserting product in db"), "error inserting product in db: %s", err.Error())
	}

	// now insert the skus for this product

	objectID := res.InsertedID.(primitive.ObjectID)

	return objectID.Hex(), nil
}

func (r repository) Update(ctx context.Context, id string, req updateRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.mdb.UpdateOne(
		"product",
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{
				{"title", req.Title},
				{"body", req.Body},
				{"category", req.Category},
				{"assets", req.Assets},
				{"variants", req.Variants}}}})
	if err != nil {
		return errors.Wrapf(errorsd.New("error updating product in db"), "error updating product in db: %s", err.Error())
	}
	return nil
}

func (r repository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.mdb.DeleteOne(
		"product",
		ctx,
		bson.M{"_id": objectID})
	if err != nil {
		return errors.Wrapf(errorsd.New("error deleting product in db"), "error deleting product in db: %s", err.Error())
	}
	return nil
}
