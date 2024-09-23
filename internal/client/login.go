package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

// LoginGet sends a get request to the login page
func (c *Client) LoginGet() io.Reader {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+LOGIN_URL, nil)
	if err != nil {
		log.Fatalf("error while logging in %v", err)
	}

	getResponse, err := c.HttpClient.Do(req)
	if err != nil {
		log.Fatalf("error while logging in %v", err)
	}

	defer getResponse.Body.Close()

	// create a buffer from body
	buffer := bytes.NewBuffer(nil)

	_, err = io.Copy(buffer, getResponse.Body)

	if err != nil {
		panic("error while copying the bufer")
	}

	return buffer
}

// Login sends a post request to log in
// We don't need to save any data after a successful login since auth will be handled with client's sessions
func (c *Client) Login(id string, password string) *bytes.Buffer {
	formData := url.Values{}
	formData.Set("txtTCPasaport", id)
	formData.Set("txtSifre", password)
	formData.Set("btnGirisYap", "Giriş Yap")
	formData.Set("__VIEWSTATE", c.ViewState)

	resp, postErr := c.HttpClient.PostForm(c.BaseURL+LOGIN_URL, formData)

	if postErr != nil {
		panic(postErr.Error())
	}

	defer resp.Body.Close()

	// create a buffer from body
	buffer := bytes.NewBuffer(nil)

	_, err := io.Copy(buffer, resp.Body)

	if err != nil {
		panic("error while copying the bufer")
	}

	return buffer
}
