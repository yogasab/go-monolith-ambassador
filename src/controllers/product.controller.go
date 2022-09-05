package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
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

func (h *productController) UpdateProduct(ctx *fiber.Ctx) error {
	dto := new(dto.UpdateProductDTO)
	ID, _ := strconv.Atoi(ctx.Params("id"))
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{
				"code":    http.StatusUnprocessableEntity,
				"message": "failed to process request",
				"error":   err,
			})
	}

	if errors := ValidateInput(*dto); errors != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": "error validation request",
			"error":   errors,
		})
	}

	dto.ID = ID

	updatedProduct, err := h.productService.UpdateProduct(dto)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "product updated successfully",
			"data":    updatedProduct,
		})
}

func (h *productController) DeleteProduct(ctx *fiber.Ctx) error {
	ID, _ := strconv.Atoi(ctx.Params("id"))

	isDeleted, err := h.productService.DeleteProduct(ID)
	if err != nil {
		if err.Error() == "product not found" {
			return ctx.Status(http.StatusNotFound).
				JSON(fiber.Map{
					"code":    http.StatusNotFound,
					"message": "not found",
					"error":   err.Error(),
				})
		}
		return ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":       http.StatusOK,
			"message":    "product deleted successfully",
			"is_deleted": isDeleted,
		})
}
