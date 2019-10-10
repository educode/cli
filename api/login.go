package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	LoginURL string = "api/auth/v1/authenticate"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(username string, password string) (string, error) {
	u, err := c.BaseURL.Parse(LoginURL)
	if err != nil {
		return "", err
	}

	credentials := Credentials{
		Username: username,
		Password: password,
	}

	body, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	auth, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authResponse AuthResponse
	if err := json.Unmarshal(auth, &authResponse); err != nil {
		return "", err
	}

	return authResponse.Token, nil
}