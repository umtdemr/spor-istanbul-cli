package client

import (
	"net/http"
	"net/url"
)

func Login(id string, password string) (string, bool) {
	formData := url.Values{}
	formData.Set("txtTCPasaport", id)
	formData.Set("txtSifre", password)
	formData.Set("btnGirisYap", "Giriş Yap")

	resp, _ := http.PostForm(LOGIN_URL, formData)

	defer resp.Body.Close()

	cookies := resp.Cookies()
	var sessionId string

	for _, cookie := range cookies {
		if cookie.Name == COOKIE_SESSION_ID {
			sessionId = cookie.Value
		}
	}

	return sessionId, sessionId != ""
}
