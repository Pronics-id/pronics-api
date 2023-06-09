package services

import (
	"context"
	"fmt"
	"os"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerService interface {
	GetCustomerProfile(ctx context.Context, ID primitive.ObjectID) (formatters.CustomerResponse, error)
	UpdateProfileCustomer(ctx context.Context, ID primitive.ObjectID, input inputs.UpdateProfilCustomerInput, fileName string) (*mongo.UpdateResult, error)
	GetAllCustomer(ctx context.Context) ([]formatters.CustomerDashboardAdminResponse, error)
}

type customerService struct {
	userRepository     repositories.UserRepository
	customerRepository repositories.CustomerRepository
	alamatCustomerRepository repositories.AlamatCustomerRepository
	orderRepository repositories.OrderRepository
}

func NewCustomerService(userRepository repositories.UserRepository, customerRepository repositories.CustomerRepository, alamatCustomerRepository repositories.AlamatCustomerRepository, orderRepository repositories.OrderRepository) *customerService{
	return &customerService{userRepository, customerRepository, alamatCustomerRepository, orderRepository}
}

func (s *customerService) GetCustomerProfile(ctx context.Context, ID primitive.ObjectID) (formatters.CustomerResponse, error){ 
	var data formatters.CustomerResponse

	user, err := s.userRepository.GetUserById(ctx, ID)

	if err != nil{
		return data, err
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, user.ID)

	if err != nil{
		return data, err
	}

	var formatAlamats []formatters.AlamatResponse

	for _, alamatId := range customer.AlamatCustomer{
		alamat, err := s.alamatCustomerRepository.GetAlamatById(ctx, alamatId)

		if err != nil{
			return data, err
		}

		alamatFormat := formatters.AlamatResponse{
			ID : alamat.ID,
			Alamat: alamat.Alamat,
			IsUtama: alamat.IsUtama,
		}

		formatAlamats = append(formatAlamats, alamatFormat)
	}

	if err != nil {
		return data, err
	}

	data = helper.MapperCustomer(user, customer, formatAlamats)

	return data, nil
}

func (s *customerService) UpdateProfileCustomer(ctx context.Context, ID primitive.ObjectID, input inputs.UpdateProfilCustomerInput, fileName string) (*mongo.UpdateResult, error){
	var newCustomer primitive.M
	
	if fileName != ""{
		newCustomer = bson.M{
			"username" : input.Username,
			"gambarcustomer": os.Getenv("CLOUD_STORAGE_READ_LINK")+"customer/"+fileName,
			"updatedat" : time.Now(),
		}
	}else{
		newCustomer = bson.M{
			"username" : input.Username,
			"updatedat" : time.Now(),
		}
	}
	

	newUser := bson.M{
		"namalengkap" : input.NamaLengkap,
		"email" : input.Email,
		"notelepon" : input.NoHandphone,
		"deskripsi" : input.Deskripsi,
		"jeniskelamin" : input.JenisKelamin,
		"tanggallahir" : input.TanggalLahir,
		"updatedat": time.Now(),
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx,ID)

	if err != nil{
		return nil, err
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, ID, newUser)

	if err != nil{
		return nil, err
	}

	updatedCustomer, err := s.customerRepository.UpdateProfil(ctx, customer.ID,newCustomer)

	if err != nil{
		return nil, err
	}

	fmt.Println(updatedCustomer)

	return updatedUser, nil
}

func (s *customerService) GetAllCustomer(ctx context.Context) ([]formatters.CustomerDashboardAdminResponse, error){
	var allCustomers []formatters.CustomerDashboardAdminResponse

	customers, err := s.customerRepository.GetAllCustomer(ctx)

	if err != nil{
		return allCustomers, err
	}

	for _, customer := range customers{
		var customerResponse formatters.CustomerDashboardAdminResponse

		user, err := s.userRepository.GetUserById(ctx, customer.UserId)

		if err != nil{
			return allCustomers, err
		}

		customerResponse.ID = customer.ID
		customerResponse.Email = user.Email
		customerResponse.NamaLengkap = user.NamaLengkap
		customerResponse.NoHandphone = user.NoTelepon

		orders, err := s.orderRepository.GetAllOrderCustomer(ctx, customer.ID)

		if err != nil{
			customerResponse.JumlahTransaksi = 0
		}else{
			customerResponse.JumlahTransaksi = len(orders)
		}

		allCustomers = append(allCustomers, customerResponse)
	}

	return allCustomers, nil
}