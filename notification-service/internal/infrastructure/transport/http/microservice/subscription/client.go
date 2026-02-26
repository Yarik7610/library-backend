package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client interface {
	GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() Client {
	return &client{
		baseURL: microservice.SUBSCRIPTIONS_HTTP_ADDRESS,
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (c *client) GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+route.SUBSCRIPTIONS+route.BOOKS+route.CATEGORIES+"/"+bookCategory, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Subscription microservice returned status code: %d", resp.StatusCode)
	}

	var emails []string
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}
	return emails, nil
}
