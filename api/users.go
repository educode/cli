package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	MembersBaseUrl string = "api/auth/v1/users"
	PageParam string = "page"
	PageSizeParam string = "pageSize"
)

type MemberInfo struct {
	Id         	string `json:"id"`
	Title      	string `json:"title"`
	Firstname	string `json:"firstname"`
	Lastname 	string `json:"lastname"`
	Email		string `json:"email"`
	Role 		string `json:"role"`
}

func (s *MemberInfo) ToRow() []string {
	return []string{s.Id, s.Title, s.Firstname, s.Lastname, s.Email, s.Role}
}

func (c *Client) GetMembers(page int, pageSize int) ([]MemberInfo, error) {
	u, err := c.BaseURL.Parse(MembersBaseUrl)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	if page != 0 { params.Add(PageParam, strconv.Itoa(page)) }
	if pageSize != 0 { params.Add(PageSizeParam, strconv.Itoa(pageSize)) }

	u.RawQuery = params.Encode()
	resp, err := c.Request(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var memberInfos []MemberInfo
	if err := json.Unmarshal(buffer, &memberInfos); err != nil {
		return nil, err
	}

	return memberInfos, err
}