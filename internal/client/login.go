package client

import (
	"github.com/umtdemr/spor-istanbul-cli/internal/parser"
	"net/http"
	"net/url"
	"strings"
)

func Login(id string, password string) (string, bool) {
	formData := url.Values{}
	formData.Set("txtTCPasaport", id)
	formData.Set("txtSifre", password)
	formData.Set("btnGirisYap", "Giriş Yap")
	formData.Set("__VIEWSTATE", "18v9/jvlC8qsN16XpBUmSb1Pq4Qp4X0pMErF1AMS0Kw/METmb6YGeh04udRG+fyrUGWFjPMGPETZp7235nCmqmDNRkAlboNzDmgy7etyxJcHXpwBY1+pxMTfnOTlOsz/")

	resp, _ := http.PostForm(LOGIN_URL, formData)

	defer resp.Body.Close()

	cookies := resp.Cookies()

	pageTitle, ok := parser.GetTitle(resp.Body)

	if !ok {
		return "", false
	}

	if strings.Contains(pageTitle, "Giriş Yap") {
		return "", false
	}

	var sessionId string

	for _, cookie := range cookies {
		if cookie.Name == COOKIE_SESSION_ID {
			sessionId = cookie.Value
		}
	}

	return sessionId, sessionId != ""
}
