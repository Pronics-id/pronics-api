package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KTPMitra struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	MitraId primitive.ObjectID `json:"mitra_id" bson:"mitra_id"`
	GambarKtp string `json:"gambar_ktp"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}