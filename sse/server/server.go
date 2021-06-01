package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
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
	r := gin.Default()

	r.GET("/stream", streamApi)
	log.Fatal(r.Run(":4040"))
}

func streamApi(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 스트림에 전송할 채널 생성
	ch := make(chan Log, 10)
	// 스트림에 전송할 데이터를 가져오는 고루틴 생성
	go sendLogging(ch)

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <- ch; ok {
			c.SSEvent("", msg)
			return true
		}
		return false
	})
}

func sendLogging(ch chan<- Log) {
	defer close(ch)
	for {
		ch <- getLog()
	}
}

// getLog 로그를 랜덤하게 반환한다
func getLog() Log{
	seed := rand.NewSource(time.Now().UnixNano())
	rn := rand.New(seed)
	size := len(Logs)

	return Logs[rn.Intn(size-1)]
}