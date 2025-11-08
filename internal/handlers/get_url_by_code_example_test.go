package handlers

import (
	"fmt"
	"net/http"
)

func ExampleGetURLByCode() {
	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := c.Get("http://localhost:8080/rC8mfESy")
	if err != nil {
		fmt.Println("failed to make request", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		fmt.Println("short url not found")
		return
	}

	if res.StatusCode != http.StatusTemporaryRedirect {
		fmt.Println("bad status", res.StatusCode)
		return
	}

	fmt.Println("location", res.Header.Get("Location"))
	fmt.Println("success")
}
