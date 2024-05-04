package client

import (
	"io"
	"net/http"
)

func (c *Client) ChooseSession() {

}

func (c *Client) GetSubscriptionsPage() io.Reader {
	req, _ := http.NewRequest(http.MethodGet, BASE_URL+SUBSCRIPTIONS_LIST_URL, nil)

	resp, _ := c.HttpClient.Do(req)

	defer resp.Body.Close()

	return resp.Body
}
