package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
)

type Client interface {
	BookCategoryExists(ctx context.Context, bookCategory string) (bool, error)
}

type client struct {
	baseURL string
}

func NewClient() Client {
	return &client{baseURL: microservice.CATALOG_ADDRESS}
}

func (c *client) BookCategoryExists(ctx context.Context, bookCategory string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+route.CATALOG+route.BOOKS+route.CATEGORIES, nil)
	if err != nil {
		return false, errs.NewInternalServerError().WithCause(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errs.NewInternalServerError().WithCause(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errs.NewInternalServerError().WithCause(fmt.Errorf("catalog microservice returned status code: %d", resp.StatusCode))
	}

	var bookCategories []string
	if err := json.NewDecoder(resp.Body).Decode(&bookCategories); err != nil {
		return false, errs.NewInternalServerError().WithCause(err)
	}

	bookCategory = strings.ToLower(bookCategory)
	return slices.Contains(bookCategories, bookCategory), nil
}
