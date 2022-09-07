package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yogasab/go-monolith-ambassador/src/middlewares"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/services"
)

type linkController struct {
	linkService services.LinkService
}

func NewLinkController(linkService services.LinkService) *linkController {
	return &linkController{linkService: linkService}
}

func (h *linkController) GetUserLinks(ctx *fiber.Ctx) error {
	UserID, _ := strconv.Atoi(ctx.Params("id"))
	links, err := h.linkService.GetUserLinks(UserID)
	if err != nil {
		if err.Error() == "user links is not found" {
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
			"code":    http.StatusOK,
			"message": "user links fetched successfully",
			"data":    links,
		})
}

func (h *linkController) CreateLink(ctx *fiber.Ctx) error {
	dto := new(dto.CreateLinkDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{
				"code":    http.StatusUnprocessableEntity,
				"message": "failed to process request",
				"error":   err,
			})
	}

	ID, _ := middlewares.GetUserID(ctx)
	dto.UserID = ID
	if errors := ValidateInput(*dto); errors != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{
				"code":    http.StatusUnprocessableEntity,
				"message": "error validation request",
				"error":   errors,
			})
	}

	newLink, err := h.linkService.CreateLink(dto)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"message": "internal server error",
				"error":   err,
			})
	}

	return ctx.
		Status(http.StatusCreated).
		JSON(fiber.Map{
			"code":    http.StatusCreated,
			"message": "user registered successfully",
			"data":    newLink,
		})
}
