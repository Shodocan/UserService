package handlers

import (
	"net/http"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/configs/engine"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/Shodocan/UserService/internal/domain/entity"
	"github.com/Shodocan/UserService/internal/injection"
	"github.com/Shodocan/UserService/internal/web/requests"
	"github.com/gofiber/fiber/v2"
)

// Search godoc
// @Summary Search Users
// @Description Search users
// @Accept  json
// @Produce  json
// @Param Request body requests.SearchUserRequest true "Search Users Request"
// @Success 200 {object} engine.PaginationResponse{data=[]entity.User}
// @Failure 400,401,404,500 {object} engine.Error
// @Router /users/search [post]
func SearchUsers(config *configs.EnvVarConfig, db database.MongoDB) fiber.Handler {
	useCase := injection.InitializeUserCase(config, db)
	return func(ctx *fiber.Ctx) error {
		var request requests.SearchUserRequest
		err := ctx.BodyParser(&request)
		if err != nil {
			return err
		}

		err = request.Validate()
		if err != nil {
			return err
		}

		users, pagination, err := useCase.Search(request.Filters, request.Sort, request.Limit, request.Page)
		if err != nil {
			return err
		}

		response := engine.NewResponsePaginated(users, pagination, "Users Found")
		return ctx.Status(http.StatusOK).JSON(response)
	}
}

// Find godoc
// @Summary Find User
// @Description Find user
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} engine.Response{data=entity.User}
// @Failure 400,401,404,500 {object} engine.Error
// @Router /users/{id} [get]
func FindUser(config *configs.EnvVarConfig, db database.MongoDB) fiber.Handler {
	useCase := injection.InitializeUserCase(config, db)
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		if id == "" {
			return engine.ErrBadRequest().Message("ID is required")
		}

		user, err := useCase.Find(id)
		if err != nil {
			return err
		}

		response := engine.NewResponseOK(user, "User Found")
		return ctx.Status(http.StatusOK).JSON(response)
	}
}

// ValidatePassword godoc
// @Summary ValidatePassword
// @Description Validate Password
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param Request body requests.ValidatePassword true "Validate Password Request"
// @Success 200 {object} engine.Response{data=string}
// @Failure 400,401,404,500 {object} engine.Error
// @Router /users/password/{id} [post]
func ValidatePassword(config *configs.EnvVarConfig, db database.MongoDB) fiber.Handler {
	useCase := injection.InitializeUserCase(config, db)
	return func(ctx *fiber.Ctx) error {
		var request requests.ValidatePassword
		err := ctx.BodyParser(&request)
		if err != nil {
			return err
		}

		if request.Password == "" {
			return engine.ErrBadRequest().Message("Password required")
		}

		id := ctx.Params("id")
		if id == "" {
			return engine.ErrBadRequest().Message("ID is required")
		}

		err = useCase.ValidatePassword(id, request.Password)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(engine.NewResponseOK("", "Valid"))
	}
}

// Create godoc
// @Summary Create User
// @Description Create user
// @Accept  json
// @Produce  json
// @Param Request body entity.User true "Create User Request"
// @Success 201 {object} engine.Response{data=entity.User}
// @Failure 400,401,404,500 {object} engine.Error
// @Router /users [post]
func CreateUser(config *configs.EnvVarConfig, db database.MongoDB) fiber.Handler {
	useCase := injection.InitializeUserCase(config, db)
	return func(ctx *fiber.Ctx) error {
		var request entity.User
		err := ctx.BodyParser(&request)
		if err != nil {
			return err
		}

		err = request.Validate(true)
		if err != nil {
			return err
		}

		user, err := useCase.Create(request)
		if err != nil {
			return err
		}
		response := engine.NewResponseCreated(user, "User Created")

		return ctx.Status(http.StatusCreated).JSON(response)
	}
}

// Update godoc
// @Summary Update User
// @Description update user
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param Request body entity.User true "Update User Request"
// @Success 200 {object} engine.Response{data=entity.User}
// @Failure 400,401,404,500 {object} engine.Error
// @Router /users/{id} [post]
func UpdateUser(config *configs.EnvVarConfig, db database.MongoDB) fiber.Handler {
	useCase := injection.InitializeUserCase(config, db)
	return func(ctx *fiber.Ctx) error {
		var request entity.User
		err := ctx.BodyParser(&request)
		if err != nil {
			return err
		}

		id := ctx.Params("id")
		if id == "" {
			return engine.ErrBadRequest().Message("ID is required")
		}

		err = request.Validate(false)
		if err != nil {
			return err
		}

		user, err := useCase.Update(id, request)
		if err != nil {
			return err
		}

		response := engine.NewResponseOK(user, "User Updated")

		return ctx.Status(http.StatusOK).JSON(response)
	}
}

// PartialUpdate godoc
// @Summary PartialUpdate User
// @Description update user attributes
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param Request body entity.User true "Update User Request"
// @Success 200 {object} engine.Response{data=entity.User}
// @Failure 400,401,404,500 {object} engine.Error
// @Router /users/{id} [put]
func PartialUpdateUser(config *configs.EnvVarConfig, db database.MongoDB) fiber.Handler {
	useCase := injection.InitializeUserCase(config, db)
	return func(ctx *fiber.Ctx) error {
		var request entity.User
		err := ctx.BodyParser(&request)
		if err != nil {
			return err
		}

		id := ctx.Params("id")
		if id == "" {
			return engine.ErrBadRequest().Message("ID is required")
		}

		user, err := useCase.PartialUpdate(id, request)
		if err != nil {
			return err
		}

		response := engine.NewResponseOK(user, "User Updated")

		return ctx.Status(http.StatusOK).JSON(response)
	}
}

// Delete godoc
// @Summary Delete User
// @Description update user attributes
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} engine.Response{data=string}
// @Failure 400,401,404,500 {object} engine.Error
// @Router /users/{id} [delete]
func Delete(config *configs.EnvVarConfig, db database.MongoDB) fiber.Handler {
	useCase := injection.InitializeUserCase(config, db)
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		if id == "" {
			return engine.ErrBadRequest().Message("ID is required")
		}

		err := useCase.Delete(id)
		if err != nil {
			return err
		}

		response := engine.NewResponseOK("", "User Deleted")

		return ctx.Status(http.StatusOK).JSON(response)
	}
}
