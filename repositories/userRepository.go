package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (models.User, error)
	Save(ctx context.Context, user models.User) (*mongo.InsertOneResult, error)
	IsUserExist(ctx context.Context, email string) (bool, error)
	GetUserById(ctx context.Context, ID primitive.ObjectID) (models.User,  error)
	UpdateUser(ctx context.Context, ID primitive.ObjectID, newUser primitive.M)(*mongo.UpdateResult, error)
}

type userRepository struct{
	DB *mongo.Collection
}

func NewUserRepository(DB *mongo.Collection) *userRepository{
	return &userRepository{DB}
}

func (r *userRepository) Save(ctx context.Context,user models.User) (*mongo.InsertOneResult, error) {
	r.DB.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys : bson.D{{Key: "email", Value: 1}},
			Options : options.Index().SetUnique(true),
		},
	)
	
	result,err := r.DB.InsertOne(ctx, user)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (models.User,  error){

	var user models.User

	err := r.DB.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil{
		return user, err
	}

	return user, nil
}

func (r *userRepository) IsUserExist(ctx context.Context, email string) (bool, error){
	var user models.User

	err := r.DB.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	
	if err != nil{
		return false, err
	}


	return true, nil
}

func (r *userRepository) GetUserById(ctx context.Context, ID primitive.ObjectID) (models.User,  error){

	var user models.User

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&user)

	if err != nil{
		return user, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, ID primitive.ObjectID, newUser primitive.M)(*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx,bson.M{"_id":ID},bson.M{"$set" : newUser})

	if err != nil{
		return result, err
	}

	return result, nil
}