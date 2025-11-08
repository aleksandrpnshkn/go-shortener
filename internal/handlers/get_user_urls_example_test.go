package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func ExampleGetUserURLs() {
	c := http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8080/api/user/urls", nil)
	if err != nil {
		fmt.Println("failed to make request", err)
	}
	req.Header.Add("Content-Type", "application/json")

	authCookie := &http.Cookie{
		Name:  "auth_token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.ZFQYhAk2o2DDE7PMJJcYHRgb74kcYvc-oSQ9J63elnQ",
	}
	req.AddCookie(authCookie)

	res, err := c.Do(req)
	if err != nil {
		fmt.Println("failed to send request", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		fmt.Println("user has no URLs")
		return
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("bad status", res.StatusCode)
		return
	}

	rawUserURLs, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("failed to read response body", err)
	}

	fmt.Println(string(rawUserURLs))
}
