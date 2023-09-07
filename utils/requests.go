package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func HtbGet(url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	jwt_token := GetHTBToken()
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "HTB-Tool")
	req.Header.Set("Host", "www.hackthebox.com")
	req.Header.Set("Authorization", "Bearer "+jwt_token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func HtbPost(url string, jsonData []byte) *http.Response {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	jwt_token := GetHTBToken()
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "HTB-Tool")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt_token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}
