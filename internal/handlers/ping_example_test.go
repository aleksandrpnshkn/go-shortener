package handlers

import (
	"fmt"
	"net/http"
)

func ExamplePingHandler() {
	c := http.Client{}
	res, err := c.Get("http://localhost:8080/ping")
	if err != nil {
		fmt.Println("failed to make request", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("bad status", res.StatusCode)
		return
	}

	fmt.Println("success")
}
