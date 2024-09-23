package client

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// GetSessions makes a request for the sessions page
func (c *Client) GetSessions(sessionPostId string) *bytes.Buffer {
	formData := url.Values{}
	formData.Set("__EVENTTARGET", sessionPostId)
	formData.Set("__VIEWSTATE", c.ViewState)

	resp, _ := c.HttpClient.PostForm(c.BaseURL+"/uyespor", formData)

	defer resp.Body.Close()

	buffer := bytes.NewBuffer(nil)
	_, err := io.Copy(buffer, resp.Body)
	if err != nil {
		panic("error while copying the buffer")
	}
	return buffer
}

// GetSubscriptionsPage makes a request for the subscriptions page
func (c *Client) GetSubscriptionsPage() *bytes.Buffer {
	req, _ := http.NewRequest(http.MethodGet, BASE_URL+SUBSCRIPTIONS_LIST_URL, nil)

	resp, _ := c.HttpClient.Do(req)

	defer resp.Body.Close()

	buffer := bytes.NewBuffer(nil)
	_, err := io.Copy(buffer, resp.Body)

	if err != nil {
		panic("error while copying the buffer")
	}

	return buffer
}
