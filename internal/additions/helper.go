package additions

import (
	"strings"

	"github.com/Reywaltz/avito_backend/internal/models/users"
)

func ValidateUser(user users.User) bool {
	if strings.TrimSpace(user.Username) == "" {
		return false
	}

	return true
}
