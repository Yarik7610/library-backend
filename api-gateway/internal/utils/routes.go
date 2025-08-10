package utils

import (
	"slices"

	"github.com/Yarik7610/library-backend/api-gateway/internal/constants"
)

func IsPrivateRoute(path string) bool {
	return slices.Contains(constants.PRIVATE_ROUTES, path)
}
