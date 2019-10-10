package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	ReviewBaseUrl  string = "api/review/v1/reviews"
)

type User struct {
	Id string `json:"id"`
}

type Review struct {
	Challenge string `json:"challengeId"`
	Reviewer User `json:"reviewer"`
	Student User `json:"student"`
	Deadline string `json:"deadline"`
	Content string `json:"content"`
	PointsRevoked bool `json:"pointsRevoked"`

}

func (c *Client) CreateReviews(reviews []Review) error {
	u, err := c.BaseURL.Parse(ReviewBaseUrl)
	if err != nil {
		return err
	}

	buffer, err := json.Marshal(reviews)
	if err != nil {
		return err
	}

	resp, err := c.Request(http.MethodPost, u.String(), bytes.NewBuffer(buffer))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("saving reviews failed")
	}

	return nil
}