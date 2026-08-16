package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/adityaonrepeat/user-management-api/internal/models"
	"github.com/adityaonrepeat/user-management-api/internal/service"
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
	res, err := h.svc.List(c.UserContext())
	if err != nil {
		return mapServiceError(err)
	}
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

func parseID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil || id < 1 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "id must be a positive integer")
	}
	return int32(id), nil
}
