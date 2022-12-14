package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/database"
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type productController struct {
	productService services.ProductService
	redisService   services.RedisService
}

func NewProductController(
	productService services.ProductService,
	redisService services.RedisService,
) *productController {
	return &productController{productService: productService, redisService: redisService}
}

var (
	keys []string = []string{"products_frontend", "products_backend"}
)

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

	// delete all the value from redis
	go database.ClearCache(keys...)
	// go routines with named function
	// go deleteCache(ctx.Context(), keys, database.RedisClient)
	// go routines with anonynouse function
	// go func(ctx context.Context, keys []string, redisClient *redis.Client) {
	// 	time.Sleep(5 * time.Second)
	// 	for _, key := range keys {
	// 		redisClient.Del(ctx, key).Err()
	// 	}
	// }(ctx.Context(), keys, database.RedisClient)

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "product updated successfully",
			"data":    updatedProduct,
		})
}

// func deleteCache(ctx context.Context, keys []string, redisClient *redis.Client) {
// 	time.Sleep(5 * time.Second)
// 	for _, key := range keys {
// 		redisClient.Del(ctx, key).Err()
// 	}
// }

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

	// delete all the value from redis
	go database.ClearCache(keys...)

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":       http.StatusOK,
			"message":    "product deleted successfully",
			"is_deleted": isDeleted,
		})
}

func (h *productController) CreateProduct(ctx *fiber.Ctx) error {
	dto := new(dto.CreateProductDTO)
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

	newProduct, err := h.productService.CreateProduct(dto)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server errors",
				"error":   err.Error(),
			})
	}

	// delete all the value from redis
	go database.ClearCache(keys...)

	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"code":    http.StatusOK,
			"message": "product updated successfully",
			"data":    newProduct,
		})
}

func (h *productController) GetProductsFrontend(ctx *fiber.Ctx) error {
	var products []*models.Product
	redisKey := "products_frontend"
	redisCtx := ctx.Context()

	results, errRedis := h.redisService.GetValue(redisCtx, redisKey)
	if errRedis != nil || results == "" {
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

		bytes, err := json.Marshal(products)
		_, errRedis = h.redisService.SetValue(redisCtx, redisKey, bytes)
		if errRedis != nil {
			return ctx.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"code":    http.StatusInternalServerError,
					"message": "failed to set from redis",
					"error":   err.Error(),
				})
		}

		return ctx.
			Status(http.StatusOK).
			JSON(fiber.Map{
				"code":    http.StatusOK,
				"message": "product fetched successfully",
				"data":    products,
			})
	}

	json.Unmarshal([]byte(results), &products)
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "product from redis fetched successfully",
		"data":    products,
	})
}

func (h *productController) GetProductsBackend(ctx *fiber.Ctx) error {
	var products []*models.Product
	redisKey := "products_backend"
	redisCtx := ctx.Context()

	results, errRedis := h.redisService.GetValue(redisCtx, redisKey)
	if errRedis != nil || results == "" {
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
		bytes, _ := json.Marshal(products)
		_, errRedis := h.redisService.SetValue(redisCtx, redisKey, bytes)
		if errRedis != nil {
			return ctx.
				Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"code":    http.StatusInternalServerError,
					"message": "failed to set from redis",
					"error":   err.Error(),
				})
		}
		return ctx.
			Status(http.StatusOK).
			JSON(fiber.Map{
				"code":    http.StatusOK,
				"message": "product fetched successfully",
				"data":    products,
			})
	}

	json.Unmarshal([]byte(results), &products)
	// if s query params is not empty string
	var searchedProducts []*models.Product
	if searchParam := ctx.Query("s"); searchParam != "" {
		loweredCaseS := strings.ToLower(searchParam)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), loweredCaseS) || strings.Contains(strings.ToLower(product.Description), loweredCaseS) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	// sort by price
	if priceParams := ctx.Query("price"); priceParams != "" {
		loweredPriceParams := strings.ToLower(priceParams)
		if loweredPriceParams == "lowest" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if loweredPriceParams == "highest" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	// pagination
	var total int = len(searchedProducts)
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	perPage := 9

	var data []*models.Product
	if total <= page*perPage && total >= (page-1)*perPage {
		// if it reaches last page
		data = searchedProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		// as long as the total is greater than page * perPage
		data = searchedProducts[(page-1)*perPage : page*perPage]
	} else {
		// if perpage is greater than anything
		data = []*models.Product{}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "product from redis fetched successfully",
		"data": fiber.Map{
			"data":      data,
			"total":     total,
			"page":      page,
			"last_page": total/perPage + 1,
		},
	})
}
