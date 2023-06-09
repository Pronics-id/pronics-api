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

type orderPaymentHandler struct {
	orderPaymentService services.OrderPaymentService
}

func NewOrderPaymentHandler(orderPaymentService services.OrderPaymentService) *orderPaymentHandler {
	return &orderPaymentHandler{orderPaymentService}
}

func (h *orderPaymentHandler) AddOrUpdateOrderPayment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderDetailId, _ := primitive.ObjectIDFromHex(c.Params("orderDetailId"))

	var input inputs.AddOrUpdateOrderPaymentInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add or update order payment failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	AddOrUpdateOrderPayment, err := h.orderPaymentService.AddOrUpdateOrderPayment(ctx, orderDetailId, input)

	if err != nil {
		response := helper.APIResponse("Add or Update order payment failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add or update order payment success", http.StatusOK, "success", AddOrUpdateOrderPayment)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *orderPaymentHandler) ConfirmPayment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderPaymentId, _ := primitive.ObjectIDFromHex(c.Params("orderPaymentId"))

	var input inputs.ConfirmPaymentInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Confirm payment failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	fileName := ""

	buktiBayar, err := c.FormFile("bukti_bayar")

	if buktiBayar != nil{
		if err != nil {
			response := helper.APIResponse("Confirm payment failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		blobFile, err := buktiBayar.Open()
	
		if err != nil {
			response := helper.APIResponse("Confirm payment failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}
	
		fileName = helper.GenerateFilename(buktiBayar.Filename)

		err = configs.StorageInit("buktiBayar").UploadFile(blobFile, fileName)
	
		if err != nil {
			response := helper.APIResponse("Confirm payment failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}
	}

	ConfirmedPayment, err := h.orderPaymentService.ConfirmPayment(ctx, orderPaymentId, input, fileName)

	if err != nil {
		response := helper.APIResponse("Confirm payment failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Confirm order payment success", http.StatusOK, "success", ConfirmedPayment)
	c.Status(http.StatusOK).JSON(response)
	return nil
}