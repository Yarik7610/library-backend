package app

import (
	"net/http"

	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/config"
	"github.com/Yarik7610/library-backend/notification-service/internal/infrastructure/observability/metrics"
)

func newHTTPServer(config *config.Config) (*http.Server, error) {
	metricsHandler, err := metrics.Init()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle(route.METRICS, metricsHandler)

	httpServer := &http.Server{
		Addr:    ":" + config.HTTPServerPort,
		Handler: mux,
	}

	return httpServer, nil
}
