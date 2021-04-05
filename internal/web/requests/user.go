package requests

import (
	"fmt"
	"net/http"

	"github.com/Shodocan/UserService/internal/configs/engine"
	"github.com/Shodocan/UserService/internal/domain/entity"
)

type ValidatePassword struct {
	Password string `json:"password"`
}

type SearchUserRequest struct {
	Filters []entity.UserFilter `json:"filters,omitempty"`
	Sort    []string            `json:"sort,omitempty" example:"-name,age"`
	Limit   int                 `json:"limit" example:"10"`
	Page    int                 `json:"page" example:"1"`
}

func (u SearchUserRequest) Validate() error {
	validationErrors := map[string]interface{}{}
	if len(u.Filters) > 0 {
		for i, filter := range u.Filters {
			if !filter.Valid() {
				validationErrors[fmt.Sprintf("filter%d", i)] = fmt.Sprintf("Invalid Filter: %s %s %s", filter.Field, filter.Operator, filter.Value)
			}
		}
	}
	if u.Limit == 0 {
		validationErrors["lmit"] = "limit is required"
	}
	if u.Page == 0 {
		validationErrors["page"] = "page is required"
	}
	if len(validationErrors) > 0 {
		return engine.NewGenericError(http.StatusBadRequest, "Invalid Request").ExtraData(validationErrors)
	}
	return nil
}
