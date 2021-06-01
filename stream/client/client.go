package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
)


type Log struct {
	Device string `json:"Device"`
	Event string `json:"Event"`
	Result string `json:"Result"`
}

func main() {
	URL := "http://localhost:4040/stream"

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("buffer read bytes failed. %s\n", err)
		}
		var msg Log
		err = json.Unmarshal(line, &msg)
		if err != nil {
			log.Fatalf("marshaling failed. %s\n", err)
		}
		log.Printf("%+v\n", msg)
	}
}