package controllers

import (
	"context"
	"net/http"
	"pronics-api/configs"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mitraHandler struct {
	mitraService services.MitraService
}

func NewMitraHandler(mitraService services.MitraService) *mitraHandler{
	return &mitraHandler{mitraService}
}

func (h *mitraHandler) GetProfile(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	mitra, err := h.mitraService.GetMitraProfile(ctx,currentUserId)

	if err != nil{
		response := helper.APIResponse("Can't get mitra profile", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("get profil mitra success", http.StatusOK, "success", mitra)
	c.Status(http.StatusOK).JSON(response)
	return nil

}

func (h *mitraHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	var input inputs.UpdateProfilMitraInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Update profil failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	fileName := ""

	profilMitra, err := c.FormFile("gambar_mitra")

	if profilMitra != nil{
		if err != nil {
			response := helper.APIResponse("Update profil failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}
	
		blobFile, err := profilMitra.Open()
	
		if err != nil {
			response := helper.APIResponse("Update profil failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}
	
		fileName = helper.GenerateFilename(profilMitra.Filename)
	
		err = configs.StorageInit("mitra").UploadFile(blobFile, fileName)
	
		if err != nil {
			response := helper.APIResponse("Update profil failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}
	}

	updatedMitra, err := h.mitraService.UpdateProfileMitra(ctx, currentUserId, input, fileName)

	if err != nil {
		response := helper.APIResponse("Update profil failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Update profil mitra success", http.StatusOK, "success", updatedMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *mitraHandler) UploadMultipleImagesToGaleri(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	form, err := c.MultipartForm()

	if err != nil{
		response := helper.APIResponse("Error upload files", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	images := form.File["galeri_image"]

	fileNames := []string{}

	for _, image := range images{
		blobFile, err := image.Open()

		if err != nil{
			response := helper.APIResponse("Error upload files", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileName := helper.GenerateFilename(image.Filename)

		err = configs.StorageInit("galeriMitra").UploadFile(blobFile, fileName)

		if err != nil{
			response := helper.APIResponse("Error upload files", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileNames = append(fileNames, fileName)
	}

	updatedGaleriMitra, err := h.mitraService.UploadGaleriImage(ctx, currentUserId, fileNames)

	if err != nil{
		response := helper.APIResponse("Error upload files", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Upload images to galeri success", http.StatusOK, "success", updatedGaleriMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *mitraHandler) GetBidangMitra(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	bidangMitra, err := h.mitraService.GetAllBidangMitra(ctx, currentUserId)

	if err != nil {
		response := helper.APIResponse("Get bidang mitra failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get bidang mitra success", http.StatusOK, "success", bidangMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *mitraHandler) UpdateBidang(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	var input inputs.UpdateBidangMitraInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Update bidang mitra failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	updatedMitra, err := h.mitraService.UpdateBidang(ctx, currentUserId, input)

	if err != nil {
		response := helper.APIResponse("Update bidang mitra failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Update bidang mitra success", http.StatusOK, "success", updatedMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// detail bidang
func (h *mitraHandler) DetailBidangMitra(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	bidangId, _ := primitive.ObjectIDFromHex(c.Params("bidangId"))

	bidangMitra, err := h.mitraService.DetailBidang(ctx, currentUserId, bidangId)

	if err != nil {
		response := helper.APIResponse("Get bidang mitra failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get bidang mitra success", http.StatusOK, "success", bidangMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
}
