package core

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ForwardTo(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		zap.S().Fatalf("Failed to parse target URL %s: %v", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		req.Header = req.Header.Clone()

		userID, ok := req.Context().Value(sharedconstants.HEADER_USER_ID).(uint)
		if ok {
			req.Header.Set(sharedconstants.HEADER_USER_ID, strconv.FormatUint(uint64(userID), 64))
		}
		isAdmin, ok := req.Context().Value(sharedconstants.HEADER_IS_ADMIN).(bool)
		if ok {
			req.Header.Set(sharedconstants.HEADER_IS_ADMIN, strconv.FormatBool(isAdmin))
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		zap.S().Errorf("API-gateway error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	return func(ctx *gin.Context) {
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
