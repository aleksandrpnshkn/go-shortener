package handlers

import (
	"bytes"
	"fmt"
	"net/http"
)

func ExampleDeleteUserURLs() {
	c := http.Client{}

	reqDataBuffer := bytes.NewBufferString(`[
		"deHb2d15",
		"twpsTBA1"
	]`)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/user/urls", reqDataBuffer)
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

	if res.StatusCode != http.StatusAccepted {
		fmt.Println("bad status", res.StatusCode)
		return
	}

	fmt.Println("success")
}
