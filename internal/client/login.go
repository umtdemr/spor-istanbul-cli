package client

import (
	"bytes"
	"io"
	"net/url"
)

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

	buffer := bytes.NewBuffer(nil)

	_, err := io.Copy(buffer, resp.Body)

	if err != nil {
		panic("error while copying the bufer")
	}

	return buffer
}
