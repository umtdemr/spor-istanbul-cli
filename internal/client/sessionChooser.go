package client

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) GetSessions(sessionPostId string) *bytes.Buffer {
	formData := url.Values{}
	formData.Set("__EVENTTARGET", sessionPostId)

	resp, _ := c.HttpClient.PostForm(c.BaseURL+"/uyespor", formData)

	defer resp.Body.Close()

	buffer := bytes.NewBuffer(nil)
	_, err := io.Copy(buffer, resp.Body)
	if err != nil {
		panic("error while copying the buffer")
	}
	return buffer
}

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
