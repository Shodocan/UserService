//+build wireinject

package injection

import (
	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/Shodocan/UserService/internal/domain/usecase"
	"github.com/Shodocan/UserService/internal/services"
	"github.com/google/wire"
)

func InitializeUserCase(config *configs.EnvVarConfig, db database.MongoDB) *usecase.UserCase {
	wire.Build(usecase.NewUserCase, configs.NewLog, services.NewUserServiceMongo)
	return &usecase.UserCase{}
}
