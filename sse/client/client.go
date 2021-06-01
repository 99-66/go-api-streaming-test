package main

import (
	"encoding/json"
	"github.com/r3labs/sse/v2"
	"log"
)

type Log struct {
	Device string `json:"Device"`
	Event string `json:"Event"`
	Result string `json:"Result"`
}

func main() {
	URL := "http://localhost:4040/stream"
	client := sse.NewClient(URL)
	//client.Connection.Transport =  &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}

	if err := client.Subscribe("log", func(msg *sse.Event) {
		var strmLog Log
		err := json.Unmarshal(msg.Data, &strmLog)
		if err != nil {
			log.Printf("unmarshaling failed. %s\n", err)
		}
		log.Printf("%+v\n", strmLog)
	}); err != nil {
		log.Fatal(err)
	}
}