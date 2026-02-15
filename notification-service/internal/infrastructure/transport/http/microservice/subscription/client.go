package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend-common/transport/http/route"
)

type Client interface {
	GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error)
}

type client struct {
	baseURL string
}

func NewClient() Client {
	return &client{baseURL: microservice.SUBSCRIPTIONS_ADDRESS}
}

func (c *client) GetBookCategorySubscribedUserEmails(ctx context.Context, bookCategory string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+route.SUBSCRIPTIONS+route.BOOKS+route.CATEGORIES+"/"+bookCategory, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
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
