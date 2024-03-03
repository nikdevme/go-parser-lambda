package auth

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func Authorization(authURL, email, password string) (*http.Client, error) {
	// Creating a cookie jar to manage cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		_ = fmt.Errorf("error creating cookie jar: %v", err)
	}

	// Creating an HTTP client with cookie support
	client := &http.Client{
		Jar: jar,
	}

	data := url.Values{}
	data.Set("EmailAddress", email)
	data.Set("Password", password)
	data.Set("RememberMe", "false")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		_ = fmt.Errorf("error creating the request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	if err != nil {
		_ = fmt.Errorf("error when executing the request: %v", err)
	}

	return client, nil
}
