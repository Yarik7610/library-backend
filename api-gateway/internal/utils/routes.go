package utils

import (
	"slices"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
)

func IsPrivateRoute(method, fullPath string) bool {
	for _, r := range sharedconstants.PRIVATE_ROUTES {
		if r.FullPath == fullPath && slices.Contains(r.AllowedMethods, method) {
			return true
		}
	}
	return false
}

func IsAdminRoute(method, fullPath string) bool {
	for _, r := range sharedconstants.ADMIN_ROUTES {
		if r.FullPath == fullPath && slices.Contains(r.AllowedMethods, method) {
			return true
		}
	}
	return false
}
