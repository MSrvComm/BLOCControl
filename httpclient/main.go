package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type endpoint struct {
	SvcName string
	Ips     []string
}

var endpoints = make(map[string][]string)

func main() {
	svcName := "testapp-svc-2"
	req, err := http.NewRequest("GET", "http://localhost:62000/"+svcName, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("error getting response: ", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error reading response: ", err.Error())
	}

	var ep endpoint
	err = json.Unmarshal(body, &ep)
	if err != nil {
		log.Fatal("error json unmarshalling: ", err.Error())
	}
	endpoints[ep.SvcName] = ep.Ips
	fmt.Println(endpoints[svcName])
}
