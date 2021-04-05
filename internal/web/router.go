package web

import (
	"net/http"

	_ "github.com/Shodocan/UserService/docs"
	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/Shodocan/UserService/internal/web/errors"
	"github.com/Shodocan/UserService/internal/web/handlers"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	keyauth "github.com/gofiber/keyauth/v2"
)

func Router(config *configs.EnvVarConfig, db database.MongoDB) *fiber.App {
	srv := fiber.New(
		fiber.Config{ErrorHandler: errors.ErrorHandler},
	)

	// simple endpoint used to check on the health status of the app
	srv.Get("/_healthz", func(ctx *fiber.Ctx) error {
		err := db.Ping()
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString("not healthy")
		} else {
			return ctx.Status(http.StatusOK).SendString("OK")
		}
	})

	srv.Get("/swagger/*", swagger.Handler) // default

	srv.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
	}))

	apiGroup := srv.Group("/api/v1")
	apiGroup.Use(logger.New())
	apiGroup.Use(cors.New(cors.Config{
		AllowOrigins:  config.AllowOrigins,
		AllowMethods:  "PUT,GET,DELETE,POST",
		AllowHeaders:  "Content-type,Authorization",
		ExposeHeaders: "Content-Length,Content-type",
		MaxAge:        36000,
	}))

	apiGroup.Use(keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, s string) (bool, error) {
			return s == config.APIToken, nil
		},
	}))

	apiRouteGroup(apiGroup, config, db)

	return srv
}

func apiRouteGroup(api fiber.Router, config *configs.EnvVarConfig, db database.MongoDB) {
	users := api.Group("/users")
	users.Post("/password/:id", handlers.ValidatePassword(config, db))
	users.Post("/search", handlers.SearchUsers(config, db))
	users.Get("/:id", handlers.FindUser(config, db))
	users.Post("/", handlers.CreateUser(config, db))
	users.Post("/:id", handlers.UpdateUser(config, db))
	users.Put("/:id", handlers.PartialUpdateUser(config, db))
	users.Delete("/:id", handlers.Delete(config, db))
}
