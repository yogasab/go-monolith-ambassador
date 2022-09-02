package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type productController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *productController {
	return &productController{productService: productService}
}

func (h *productController) GetProducts(ctx *fiber.Ctx) error {
	products, err := h.productService.GetProducts()
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "products fetched successfully",
		"data":    products,
	})
}
