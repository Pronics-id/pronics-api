package routes

import (
	"context"
	"net/http"
	"os"
	"pronics-api/configs"
	"pronics-api/controllers"
	"pronics-api/repositories"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoute(api fiber.Router, userCollection *mongo.Collection, customerCollection *mongo.Collection, mitraCollection *mongo.Collection, rekeningCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	rekeningRepository := repositories.NewRekeningRepository(rekeningCollection)

	// services
	userService := services.NewUserService(userRepository, customerRepository, mitraRepository, rekeningRepository)

	// controllers
	userHandler := controllers.NewUserHandler(userService)

	authUser := api.Group("/auth/user")

	// login and register
	authUser.Post("/register", userHandler.Register)
	authUser.Post("/login", userHandler.Login)
	authUser.Post("/registerMitra", userHandler.RegisterMitra)
	authUser.Get("/login/google",func(c *fiber.Ctx) error {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()
		url := configs.GoogleOAuthConfig().AuthCodeURL(os.Getenv("oAuth_String"))
		
		c.Redirect(url, http.StatusTemporaryRedirect)
		return nil
	})

	authUser.Get("/callback",userHandler.Callback)

}