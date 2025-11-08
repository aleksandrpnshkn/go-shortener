package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ExampleCreateShortURL() {
	c := http.Client{}

	reqDataBuffer := bytes.NewBufferString(`{
		"url": "http://example.com"
	}`)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/shorten", reqDataBuffer)
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

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusConflict {
		fmt.Println("bad status", res.StatusCode)
		return
	}

	resDataRaw, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("failed to read response body", err)
	}

	fmt.Println("http status", res.StatusCode)
	fmt.Println(string(resDataRaw))
}
