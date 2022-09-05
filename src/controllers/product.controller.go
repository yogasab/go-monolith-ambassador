package controllers

import (
	"net/http"
	"strconv"

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

func (h *productController) GetProduct(ctx *fiber.Ctx) error {
	ID, _ := strconv.Atoi(ctx.Params("id"))
	product, err := h.productService.GetProduct(ID)
	if err != nil {
		if err.Error() == "product not found" {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"code":    http.StatusNotFound,
				"message": "not found",
				"error":   err.Error(),
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    http.StatusInternalServerError,
			"message": "internal server errors",
			"error":   err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "product fetched successfully",
		"data":    product,
	})
}
