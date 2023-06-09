package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RekeningRoute(api fiber.Router, rekeningCollection *mongo.Collection, bankCollection *mongo.Collection, mitraCollection *mongo.Collection) {
	// repositories
	rekeningRepository := repositories.NewRekeningRepository(rekeningCollection)
	bankRepository := repositories.NewBankRepository(bankCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)

	// services
	rekeningService := services.NewRekeningService(rekeningRepository, bankRepository, mitraRepository)

	// controllers
	rekeningHandler := controllers.NewRekeningHandler(rekeningService)

	// auth
	allAuth := middlewares.AuthAll()

	rekening := api.Group("/rekening")

	rekening.Get("/detail", allAuth.AuthAll, rekeningHandler.GetDetailRekening)
	rekening.Put("/update", allAuth.AuthAll, rekeningHandler.ChangeDetailRekening)
	rekening.Post("/Save", allAuth.AuthAll, rekeningHandler.AddRekening)
	rekening.Get("/mitra/:mitraId", allAuth.AuthAll, rekeningHandler.GetRekeningMitra)
}