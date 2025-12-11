package main

import (
	"errors"
	"net/url"
	"url-shortener/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type App struct {
	service *service.URLShortener
}

func main() {
	svc, err := service.New()
	if err != nil {
		panic("Ошибка инициализации: " + err.Error())
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	server := &App{service: svc}

	api := app.Group("/api/v1")

	// Links
	api.Get("/links", server.GetLinks)
	api.Get("/links/:short_code", server.GetLink)
	api.Post("/links", server.CreateLink)
	api.Put("/links/:short_code", server.UpdateLink)
	api.Delete("/links/:short_code", server.DeleteLink)

	// Users
	api.Get("/users", server.GetUsers)
	api.Post("/users", server.CreateUser)

	println("🚀 Сервер запущен на http://localhost:3000")
	app.Listen(":3000")
}

func (a *App) GetLinks(c *fiber.Ctx) error {
	return c.JSON(a.service.Links())
}

func (a *App) GetLink(c *fiber.Ctx) error {
	code := c.Params("short_code")
	link, err := a.service.GetLinkByShortCode(code)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Ссылка не найдена"})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(link)
}

type CreateLinkReq struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	UserID      string `json:"user_id"`
}

func (a *App) CreateLink(c *fiber.Ctx) error {
	var req CreateLinkReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Некорректный JSON"})
	}

	if req.OriginalURL == "" || req.ShortCode == "" {
		return c.Status(400).JSON(fiber.Map{"error": "original_url и short_code обязательны"})
	}

	if !isValidURL(req.OriginalURL) {
		return c.Status(400).JSON(fiber.Map{
			"error": "original_url должен быть корректным URL (например, https://example.com)",
		})
	}

	userID := req.UserID
	if userID == "" {
		userID = "anonymous"
	}

	err := a.service.AddLink(req.OriginalURL, req.ShortCode, userID)
	if err != nil {
		if errors.Is(err, service.ErrShortCodeExists) {
			return c.Status(400).JSON(fiber.Map{"error": "Короткий код уже занят"})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	link, _ := a.service.GetLinkByShortCode(req.ShortCode)
	return c.Status(201).JSON(link)
}

type UpdateLinkReq struct {
	OriginalURL string `json:"original_url"`
}

func (a *App) UpdateLink(c *fiber.Ctx) error {
	code := c.Params("short_code")
	var req UpdateLinkReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Некорректный JSON"})
	}

	if req.OriginalURL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "original_url обязателен"})
	}

	if !isValidURL(req.OriginalURL) {
		return c.Status(400).JSON(fiber.Map{"error": "Некорректный URL"})
	}

	err := a.service.UpdateLink(code, req.OriginalURL)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Ссылка не найдена"})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	link, _ := a.service.GetLinkByShortCode(code)
	return c.JSON(link)
}

func (a *App) DeleteLink(c *fiber.Ctx) error {
	code := c.Params("short_code")
	err := a.service.DeleteLink(code)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Ссылка не найдена"})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func (a *App) GetUsers(c *fiber.Ctx) error {
	return c.JSON(a.service.Users())
}

type CreateUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (a *App) CreateUser(c *fiber.Ctx) error {
	var req CreateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Некорректный JSON"})
	}

	if req.Name == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name и email обязательны"})
	}

	err := a.service.AddUser(req.Name, req.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserEmailExists) {
			return c.Status(400).JSON(fiber.Map{"error": "Пользователь с таким email уже существует"})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	user, _ := a.service.GetUserByEmail(req.Email)
	return c.Status(201).JSON(user)
}

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}
