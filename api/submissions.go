package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	SubmissionBaseUrl string = "api/submission/v1/submissions"
	UserParam         string = "user"
	ChallengeParam    string = "challenge"
)

type SubmissionInfo struct {
	Challenge string `json:"challenge"`
	User string `json:"user"`
	Points int `json:"points"`
}

func (s *SubmissionInfo) ToRow() []string {
	return []string{s.Challenge, s.User, strconv.Itoa(s.Points)}
}

func (c *Client) GetSubmissions(user string, challenge string) ([]SubmissionInfo, error) {
	u, err := c.BaseURL.Parse(SubmissionBaseUrl)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	if len(user) > 0 { params.Add(UserParam, user) }
	if len(challenge) > 0 { params.Add(ChallengeParam, challenge) }

	u.RawQuery = params.Encode()
	resp, err := c.Request(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fetching submissions failed")
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var submissions []SubmissionInfo
	if err := json.Unmarshal(buffer, &submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}