package utils

import (
	"slices"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
)

func IsPrivateRoute(path string) bool {
	return slices.Contains(sharedconstants.PRIVATE_ROUTES, path)
}
