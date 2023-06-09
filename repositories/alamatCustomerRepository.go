package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlamatCustomerRepository interface {
	Save(ctx context.Context, alamat models.AlamatCustomer) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.AlamatCustomer, error)
	FindAllByCustomerId(ctx context.Context, customerId primitive.ObjectID) ([]models.AlamatCustomer, error)
	GetAlamatById(ctx context.Context, ID primitive.ObjectID) (models.AlamatCustomer,  error)
	GetAlamatUtamaCustomer(ctx context.Context, ID primitive.ObjectID) (models.AlamatCustomer,  error)
	UpdateAlamat(ctx context.Context, IdAlamat primitive.ObjectID, newAlamat primitive.M) (*mongo.UpdateResult, error)
	DeleteAlamat(ctx context.Context, IdAlamat primitive.ObjectID) (*mongo.DeleteResult, error)
}

type alamatCustomerRepository struct {
	DB *mongo.Collection
}

func NewAlamatCustomerRepository(DB *mongo.Collection) *alamatCustomerRepository {
	return &alamatCustomerRepository{DB}
}

// get all alamat
func (r *alamatCustomerRepository) FindAll(ctx context.Context) ([]models.AlamatCustomer, error){
	var alamats []models.AlamatCustomer

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return alamats, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var alamat models.AlamatCustomer
        err := currentRes.Decode(&alamat)
        if err != nil {
            return alamats, err
        }

        alamats =append(alamats, alamat)
    }

	if err := currentRes.Err(); err != nil {
        return alamats, err
    }

	currentRes.Close(ctx)

	return alamats, nil
}

// find all by customer id
func (r *alamatCustomerRepository) FindAllByCustomerId(ctx context.Context, customerId primitive.ObjectID) ([]models.AlamatCustomer, error){
	var alamats []models.AlamatCustomer

	currentRes, err := r.DB.Find(ctx, bson.M{"customer_id": customerId})

	if err != nil{
		return alamats, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var alamat models.AlamatCustomer
        err := currentRes.Decode(&alamat)
        if err != nil {
            return alamats, err
        }

        alamats =append(alamats, alamat)
    }

	if err := currentRes.Err(); err != nil {
        return alamats, err
    }

	currentRes.Close(ctx)

	return alamats, nil
}

// add alamat
func (r *alamatCustomerRepository) Save(ctx context.Context, alamat models.AlamatCustomer) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, alamat)

	if err != nil {
		return result, err
	}

	return result, nil
}

// get alamat by id
func (r *alamatCustomerRepository) GetAlamatById(ctx context.Context, ID primitive.ObjectID) (models.AlamatCustomer,  error){

	var alamat models.AlamatCustomer

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&alamat)

	if err != nil{
		return alamat, err
	}

	return alamat, nil
}

// get alamat utama
func (r *alamatCustomerRepository) GetAlamatUtamaCustomer(ctx context.Context, customerId primitive.ObjectID) (models.AlamatCustomer,  error){

	var alamat models.AlamatCustomer

	err := r.DB.FindOne(ctx, bson.M{"customer_id":customerId,"isutama": true, }).Decode(&alamat)

	if err != nil{
		return alamat, err
	}

	return alamat, nil
}


// edit alamat status utama
func (r *alamatCustomerRepository) UpdateAlamat(ctx context.Context, IdAlamat primitive.ObjectID, newAlamat primitive.M) (*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx,bson.M{"_id":IdAlamat},bson.M{"$set" : newAlamat})

	if err != nil{
		return result, err
	}

	return result, nil
}

// delete alamat
func (r *alamatCustomerRepository) DeleteAlamat(ctx context.Context, IdAlamat primitive.ObjectID) (*mongo.DeleteResult, error){
	result, err := r.DB.DeleteOne(ctx,bson.M{"_id":IdAlamat})

	if err != nil{
		return result, err
	}

	return result, nil
}