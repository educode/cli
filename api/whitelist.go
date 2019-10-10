package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	WhitelistBaseUrl string = "api/auth/v1/whitelist"
)

type Members struct {
	Usernames []string `json:"usernames"`
}

func (c *Client) SyncWhitelist(members Members) error {
	u, err := c.BaseURL.Parse(WhitelistBaseUrl)
	if err != nil {
		return err
	}

	buffer, err := json.Marshal(members)
	if err != nil {
		return err
	}

	resp, err := c.Request(http.MethodPut, u.String(), bytes.NewBuffer(buffer))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("syncing whitelist failed")
	}

	return nil;
}