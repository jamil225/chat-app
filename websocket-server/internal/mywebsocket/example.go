package mywebsocket

import (
	"log"
	"websocket/internal/httpclient"
)

func CallExternalAPI() {
	client := httpclient.NewHttpClient(10 * 1000) // 10 seconds timeout
	url := "https://jsonplaceholder.typicode.com/posts"

	// Example GET request
	resp, err := client.Get(url, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	err = httpclient.ParseResponseBody(resp, &result)
	if err != nil {
		log.Fatalf("Failed to parse response body: %v", err)
	}
	log.Printf("GET Response: %v", result)

	// Example POST request
	postData := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}
	resp, err = client.Post(url, postData, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		log.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	var postResult map[string]interface{}
	err = httpclient.ParseResponseBody(resp, &postResult)
	if err != nil {
		log.Fatalf("Failed to parse response body: %v", err)
	}
	log.Printf("POST Response: %v", postResult)
}
