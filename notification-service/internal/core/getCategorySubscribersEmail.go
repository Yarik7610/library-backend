package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Yarik7610/library-backend-common/sharedconstants"
)

func (n *notificator) getCategorySubscribersEmails(category string) ([]string, error) {
	resp, err := http.Get(sharedconstants.SUBSCRIPTIONS_MICROSERVICE_SOCKET + sharedconstants.SUBSCRIPTIONS_ROUTE + sharedconstants.CATEGORIES_ROUTE + "/" + category)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("subscription microservice return status code: %d", resp.StatusCode)
	}

	var emails []string
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}
	return emails, nil
}
