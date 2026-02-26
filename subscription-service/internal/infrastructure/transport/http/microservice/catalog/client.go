package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client interface {
	BookCategoryExists(ctx context.Context, bookCategory string) (bool, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() Client {
	return &client{
		baseURL: microservice.CATALOG_HTTP_ADDRESS,
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (c *client) BookCategoryExists(ctx context.Context, bookCategory string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+route.CATALOG+route.BOOKS+route.CATEGORIES+"/exists"+"/"+bookCategory, nil)
	if err != nil {
		return false, errs.NewInternalServerError().WithCause(err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, errs.NewInternalServerError().WithCause(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errs.NewInternalServerError().WithCause(fmt.Errorf("Catalog microservice returned status code: %d", resp.StatusCode))
	}

	var exists bool
	if err := json.NewDecoder(resp.Body).Decode(&exists); err != nil {
		return false, errs.NewInternalServerError().WithCause(err)
	}

	return exists, nil
}
