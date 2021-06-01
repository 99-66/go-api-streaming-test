package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)


type Log map[string]string
var Logs = []Log{
	{"Event": "request", "Device": "00:11:51:AA", "Result": "accepted"},
	{"Event": "response", "Device": "00:11:51:AA", "Result": "served"},
	{"Event": "response", "Device": "00:11:22:33", "Result": "served"},
	{"Event": "request", "Device": "00:11:22:33", "Result": "accepted"},
}

func main() {
	http.HandleFunc("/stream", func(w http.ResponseWriter, req *http.Request) {
		// Context cancellation
		ctx := req.Context()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		f, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")

		// 스트림에 전송할 채널 생성
		ch := make(chan Log, 10)

		go sendLogging(ctx, ch)

		for v := range ch {
			b, err := json.Marshal(v)
			if err != nil {
				log.Printf("could not json marshall reponse item %#v: %v\n", v, err)
				continue
			}

			fmt.Fprintf(w, "%s\n", string(b))
			f.Flush()
		}
	})

	log.Println("Server listening on: 4040")
	log.Fatal(http.ListenAndServe(":4040", nil))
}


func sendLogging(ctx context.Context, ch chan<- Log) {
	defer close(ch)
	for {
		select {
		// Listening for the cancellation event
		case <-ctx.Done():
			log.Printf("close: stop goroutine")
			return
		// default send to channel
		default:
			ch <- getLog()
		}
	}
}

// getLog 로그를 랜덤하게 반환한다
func getLog() Log{
	seed := rand.NewSource(time.Now().UnixNano())
	rn := rand.New(seed)
	size := len(Logs)

	return Logs[rn.Intn(size-1)]
}