package client

import (
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	BaseURL    string
	HttpClient *http.Client
	ViewState  string
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		BaseURL: BASE_URL,
		HttpClient: &http.Client{
			Jar: jar,
		},
	}
}
