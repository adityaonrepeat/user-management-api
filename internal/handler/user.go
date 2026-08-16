package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/adityaonrepeat/user-management-api/internal/models"
	"github.com/adityaonrepeat/user-management-api/internal/service"
)

const (
	totalCountHeader = "X-Total-Count"
	defaultLimit     = 10
	maxLimit         = 100
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "request body must be valid JSON")
	}
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, validationMessage(err))
	}

	res, err := h.svc.Create(c.UserContext(), req)
	if err != nil {
		return mapServiceError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	res, err := h.svc.GetByID(c.UserContext(), id)
	if err != nil {
		return mapServiceError(err)
	}
	return c.JSON(res)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	params, err := parseListParams(c)
	if err != nil {
		return err
	}

	res, total, err := h.svc.List(c.UserContext(), params)
	if err != nil {
		return mapServiceError(err)
	}

	c.Set(totalCountHeader, strconv.FormatInt(total, 10))
	return c.JSON(res)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "request body must be valid JSON")
	}
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, validationMessage(err))
	}

	res, err := h.svc.Update(c.UserContext(), id, req)
	if err != nil {
		return mapServiceError(err)
	}
	return c.JSON(res)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	if err := h.svc.Delete(c.UserContext(), id); err != nil {
		return mapServiceError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseListParams(c *fiber.Ctx) (models.ListParams, error) {
	rawLimit, rawOffset := c.Query("limit"), c.Query("offset")

	if rawLimit == "" && rawOffset == "" {
		return models.ListParams{}, nil
	}

	params := models.ListParams{Paginated: true, Limit: defaultLimit}

	if rawLimit != "" {
		v, err := strconv.ParseInt(rawLimit, 10, 32)
		if err != nil || v < 1 || v > maxLimit {
			return models.ListParams{}, fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("limit must be an integer between 1 and %d", maxLimit))
		}
		params.Limit = int32(v)
	}

	if rawOffset != "" {
		v, err := strconv.ParseInt(rawOffset, 10, 32)
		if err != nil || v < 0 {
			return models.ListParams{}, fiber.NewError(fiber.StatusBadRequest,
				"offset must be a non-negative integer")
		}
		params.Offset = int32(v)
	}

	return params, nil
}

func parseID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil || id < 1 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "id must be a positive integer")
	}
	return int32(id), nil
}
