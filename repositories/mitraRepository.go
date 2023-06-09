package repositories

import (
	"context"
	"pronics-api/constants"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MitraRepository interface {
	SaveMitra(ctx context.Context, mitra models.Mitra) (*mongo.InsertOneResult, error)
	GetMitraById(ctx context.Context, ID primitive.ObjectID) (models.Mitra,  error)
	GetMitraByIdUser(ctx context.Context, IdUser primitive.ObjectID) (models.Mitra,  error)
	UpdateProfil(ctx context.Context, ID primitive.ObjectID, newMitra primitive.M)(*mongo.UpdateResult, error)
	FindAllActiveMitra(ctx context.Context) ([]models.Mitra, error)
	GetAllMitra(ctx context.Context) ([]models.Mitra, error)
}

type mitraRepository struct{
	DB *mongo.Collection
}

func NewMitraRepository(DB *mongo.Collection) *mitraRepository{
	return &mitraRepository{DB}
}

func (r *mitraRepository) SaveMitra(ctx context.Context, mitra models.Mitra) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, mitra)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *mitraRepository) GetMitraById(ctx context.Context, ID primitive.ObjectID) (models.Mitra,  error){

	var mitra models.Mitra

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&mitra)

	if err != nil{
		return mitra, err
	}

	return mitra, nil
}

func (r *mitraRepository) GetMitraByIdUser(ctx context.Context, IdUser primitive.ObjectID) (models.Mitra,  error){

	var mitra models.Mitra

	err := r.DB.FindOne(ctx, bson.M{"user_id": IdUser}).Decode(&mitra)

	if err != nil{
		return mitra, err
	}

	return mitra, nil
}

func (r *mitraRepository) UpdateProfil(ctx context.Context, ID primitive.ObjectID, newMitra primitive.M)(*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx,bson.M{"_id":ID},bson.M{"$set" : newMitra})

	if err != nil{
		return result, err
	}

	return result, nil
}

// get all active mitra
func (r *mitraRepository) FindAllActiveMitra(ctx context.Context) ([]models.Mitra, error){
	var katalogMitras []models.Mitra
	
	currentRes, err := r.DB.Find(ctx, bson.M{"status": constants.MitraActive})

	if err != nil{
		return nil, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var Mitra models.Mitra
        err := currentRes.Decode(&Mitra)
        if err != nil {
            return katalogMitras, err
        }

        katalogMitras =append(katalogMitras, Mitra)
    }

	if err := currentRes.Err(); err != nil {
        return katalogMitras, err
    }

	currentRes.Close(ctx)

	return katalogMitras, nil
}

func (r *mitraRepository) GetAllMitra(ctx context.Context) ([]models.Mitra, error){
	var mitras []models.Mitra

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return mitras, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var mitra models.Mitra
        err := currentRes.Decode(&mitra)
        if err != nil {
            return mitras, err
        }

        mitras =append(mitras, mitra)
    }

	if err := currentRes.Err(); err != nil {
        return mitras, err
    }

	currentRes.Close(ctx)

	return mitras, nil
}

