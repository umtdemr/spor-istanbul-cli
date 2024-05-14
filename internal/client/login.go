package client

import (
	"bytes"
	"io"
	"net/url"
)

// Login sends a post request to log in
// We don't need to save any data after a successful login since auth will be handled with client's sessions
func (c *Client) Login(id string, password string) *bytes.Buffer {
	formData := url.Values{}
	formData.Set("txtTCPasaport", id)
	formData.Set("txtSifre", password)
	formData.Set("btnGirisYap", "Giriş Yap")
	formData.Set("__VIEWSTATE", "18v9/jvlC8qsN16XpBUmSb1Pq4Qp4X0pMErF1AMS0Kw/METmb6YGeh04udRG+fyrUGWFjPMGPETZp7235nCmqmDNRkAlboNzDmgy7etyxJcHXpwBY1+pxMTfnOTlOsz/")

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
