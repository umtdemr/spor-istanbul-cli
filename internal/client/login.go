package client

import (
	"github.com/umtdemr/spor-istanbul-cli/internal/parser"
	"net/url"
	"strings"
)

func (c *Client) Login(id string, password string) bool {
	formData := url.Values{}
	formData.Set("txtTCPasaport", id)
	formData.Set("txtSifre", password)
	formData.Set("btnGirisYap", "Giriş Yap")
	formData.Set("__VIEWSTATE", "18v9/jvlC8qsN16XpBUmSb1Pq4Qp4X0pMErF1AMS0Kw/METmb6YGeh04udRG+fyrUGWFjPMGPETZp7235nCmqmDNRkAlboNzDmgy7etyxJcHXpwBY1+pxMTfnOTlOsz/")

	resp, _ := c.HttpClient.PostForm(c.BaseURL+LOGIN_URL, formData)

	defer resp.Body.Close()

	pageTitle, ok := parser.GetTitle(resp.Body)

	if !ok {
		return false
	}

	if strings.Contains(pageTitle, "Giriş Yap") {
		return false
	}

	return true
}
