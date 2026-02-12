package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Yarik7610/library-backend-common/microservice"
	"github.com/Yarik7610/library-backend-common/transport/http/route"
	"github.com/Yarik7610/library-backend/subscription-service/internal/infrastructure/errs"
)

type Client interface {
	GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error)
}

type client struct {
	baseURL string
}

func NewClient() Client {
	return &client{baseURL: microservice.USER_ADDRESS}
}

func (c *client) GetEmailsByUserIDs(ctx context.Context, userIDs []uint) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+route.EMAILS, nil)
	if err != nil {
		return nil, errs.NewInternalServerError().WithCause(err)
	}

	query := req.URL.Query()
	for _, userID := range userIDs {
		query.Add("ids", strconv.FormatUint(uint64(userID), 10))
	}
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.NewInternalServerError().WithCause(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.NewInternalServerError().WithCause(fmt.Errorf("user microservice returned status code: %d", resp.StatusCode))
	}

	var emails []string
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, errs.NewInternalServerError().WithCause(err)
	}
	return emails, nil
}
