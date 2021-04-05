package errors

import (
	"github.com/Shodocan/UserService/internal/configs/engine"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	switch domainErr := err.(type) {
	case *engine.Error:
		return ctx.Status(domainErr.Code).JSON(domainErr)
	case engine.Error:
		return ctx.Status(domainErr.Code).JSON(domainErr)
	default:
		newError := engine.NewGenericError(500, "Internal Error")
		return ctx.Status(newError.Code).JSON(newError)
	}
}
