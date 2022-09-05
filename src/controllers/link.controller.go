package controllers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
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
