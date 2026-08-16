package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/adityaonrepeat/user-management-api/internal/handler"
)

func Register(app *fiber.App, users *handler.UserHandler) {
	app.Post("/users", users.Create)
	app.Get("/users", users.List)
	app.Get("/users/:id", users.GetByID)
	app.Put("/users/:id", users.Update)
	app.Delete("/users/:id", users.Delete)
}
