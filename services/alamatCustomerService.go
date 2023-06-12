package services

import (
	"context"
	"fmt"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlamatCustomerService interface {
	SaveAlamat(ctx context.Context, alamat inputs.AddAlamatCustomerInput, userId primitive.ObjectID) (*mongo.InsertOneResult, error)
	GetAllAlamat(ctx context.Context, ID primitive.ObjectID) ([]models.AlamatCustomer, error)
}

type alamatCustomerService struct{
	alamatCustomerRepository repositories.AlamatCustomerRepository
	customerRepository repositories.CustomerRepository
	userRepository     repositories.UserRepository
}

func NewAlamatCustomerService(alamatCustomerRepository repositories.AlamatCustomerRepository, customerRepository repositories.CustomerRepository, userRepository repositories.UserRepository) *alamatCustomerService{
	return &alamatCustomerService{alamatCustomerRepository, customerRepository, userRepository}
}

func (s *alamatCustomerService) SaveAlamat(ctx context.Context, alamat inputs.AddAlamatCustomerInput, userId primitive.ObjectID) (*mongo.InsertOneResult, error){

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

	if err != nil{
		return nil, err
	}
	allAlamat, err := s.alamatCustomerRepository.FindAllByCustomerId(ctx, customer.ID)

	if err != nil{
		return nil, err
	}

	isUtama := false

	if len(allAlamat) == 0{
		isUtama = true
	}
	newAlamat := models.AlamatCustomer{
		ID : primitive.NewObjectID(),
		CustomerId: customer.ID,
		Alamat : alamat.Alamat,
		IsUtama: isUtama,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	alamatAdded, err := s.alamatCustomerRepository.Save(ctx, newAlamat)

	if err != nil{
		return nil, err
	}

	var alamatArr []primitive.ObjectID

	if customer.AlamatCustomer != nil{
		alamatArr = append(alamatArr, customer.AlamatCustomer...)
	}


	alamatArr = append(alamatArr, newAlamat.ID)

	newAlamatInCustomer := bson.M{
		"alamatcustomer" : alamatArr,
		"updatedat" : time.Now(),
	}

	insertedAlamat, err := s.customerRepository.UpdateAlamatCustomer(ctx, customer.ID, newAlamatInCustomer)

	fmt.Println(insertedAlamat, err)

	return alamatAdded, nil
}

func (s *alamatCustomerService) GetAllAlamat(ctx context.Context, ID primitive.ObjectID) ([]models.AlamatCustomer, error){
	var data []models.AlamatCustomer

	user, err := s.userRepository.GetUserById(ctx, ID)

	if err != nil{
		return data, err
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, user.ID)

	if err != nil{
		return data, err
	}

	alamats, err := s.alamatCustomerRepository.FindAllByCustomerId(ctx, customer.ID)

	if err != nil{
		return data, err
	}

	return alamats, nil
}