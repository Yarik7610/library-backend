package core

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/Yarik7610/library-backend-common/transport/http/header"
	"github.com/Yarik7610/library-backend/api-gateway/internal/infrastructure/observability/logging"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func ForwardTo(logger *logging.Logger, target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		logger.Fatal(context.Background(),
			"Forward URL parse error",
			logging.String("target", target),
			logging.Error(err),
		)
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		req.Header = req.Header.Clone()

		// Cast to carrier type
		carrier := propagation.HeaderCarrier(req.Header)
		// Enrich carrier with current ctx
		otel.GetTextMapPropagator().Inject(req.Context(), carrier)

		userID, ok := req.Context().Value(header.USER_ID).(uint)
		if ok {
			req.Header.Set(header.USER_ID, strconv.FormatUint(uint64(userID), 64))
		}
		isAdmin, ok := req.Context().Value(header.IS_ADMIN).(bool)
		if ok {
			req.Header.Set(header.IS_ADMIN, strconv.FormatBool(isAdmin))
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		ctx := r.Context()

		logger.Error(ctx,
			"API-gateway error",
			logging.String("target", target),
			logging.Error(err),
		)

		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
