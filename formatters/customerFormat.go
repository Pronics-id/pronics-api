package formatters

import (
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerResponse struct{
	ID primitive.ObjectID `json:"id"`
	User UserResponse `json:"user_data"`
	Username string `json:"username"`
	GambarUser string `json:"gambar_user"`
	Alamat []AlamatResponse `json:"alamat"`
}

type CustomerDashboardAdminResponse struct{
	ID primitive.ObjectID `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email string `json:"email"`
	NoHandphone string `json:"no_handphone"`
	JumlahTransaksi int `json:"jumlah_transaksi"`
}

type CustomerFormatter struct {
	ID          primitive.ObjectID `json:"id"`
	Username  string `json:"username"`
	Email       string `json:"email"`
}

func FormatCustomer(customer models.Customer) CustomerFormatter {
	formatter := CustomerFormatter{
		ID : customer.ID,
		Username:   customer.Username,
	}

	return formatter
}
