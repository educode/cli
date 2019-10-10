package api

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"syscall"
)

const (
	baseUrl        string = "https://educode.cs.hhu.de/"
	AuthorizationHeader = "Authorization"
	ContentTypeHeader = "Content-Type"
)

type Client struct {
	BaseURL *url.URL
	httpClient *http.Client
	Token string
}

func NewClient(client *http.Client) (*Client, error) {
	// Create a default client if none is specified
	if client == nil {
		client = http.DefaultClient
	}

	// Parse the base url and create the http client
	base, _ := url.Parse(baseUrl)
	ret := &Client{httpClient: client, BaseURL: base}

	token, present := os.LookupEnv("EDUCODE_TOKEN")
	if present {
		ret.Token = "Bearer " + token
		return ret, nil
	}

	// Get the username
	user, present := os.LookupEnv("EDUCODE_USER")
	if present == false {
		fmt.Fprint(os.Stderr, "Please specify your educode username: ")
		inputBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		user = string(inputBytes)
		fmt.Fprintln(os.Stderr)
	}

	// Get the password
	password, present := os.LookupEnv("EDUCODE_PASS")
	if present == false {
		fmt.Fprint(os.Stderr, "Please specify your educode password: ")
		inputBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		password = string(inputBytes)
		fmt.Fprintln(os.Stderr)
	}

	// Login using the client
	token, err := ret.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}

	ret.Token = "Bearer " + token

	return ret, nil
}

func (c *Client) Request(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(ContentTypeHeader, "application/json")
	req.Header.Add(AuthorizationHeader, c.Token)

	return c.httpClient.Do(req)
}
