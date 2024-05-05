package client

import (
	"bytes"
	"io"
	"net/http"
)

func (c *Client) ChooseSession() {

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
