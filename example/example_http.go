package main

import (
	"encoding/json"
	"github.com/vladimirvivien/automi/emitters"
	"github.com/vladimirvivien/automi/stream"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Log map[string]string
var Logs = []map[string]string{
	{"Event": "request", "Device": "00:11:51:AA", "Result": "accepted"},
	{"Event": "response", "Device": "00:11:51:AA", "Result": "served"},
	{"Event": "response", "Device": "00:11:22:33", "Result": "served"},
	{"Event": "request", "Device": "00:11:22:33", "Result": "accepted"},
}

func main() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Content-Type", "text/html")
		resp.WriteHeader(http.StatusOK)

		// 스트림에 전송할 채널 생성
		ch := make(chan Log)
		// 스트림 에러 발생시에 고루틴을 중지하기 위한 채널 생성
		errNotify := make(chan struct{})
		// Connection 종료시에 고루틴을 중지하기 위한 채널 생성
		closeNotify := resp.(http.CloseNotifier).CloseNotify()
		go sendLogging(ch, errNotify, closeNotify)

		// Channel을 전달하는 스트림 생성
		strm := stream.New(emitters.Chan(ch))

		// 스트림에 데이터 전송을 위한 프로세스 처리
		// Into를 통해 http 응답으로 전송
		strm.Process(func(data Log) []byte {
			jsonData, _ := json.Marshal(data)
			return jsonData
		}).Into(resp)

		//strm.Into(collectors.Func(func(data interface{}) error {
		//	buf := data.([]byte)
		//	output := bytes.NewBuffer(buf)
		//	_, ok := resp.Write(append(output.Bytes(), []byte("\r")...))
		//	if ok != nil {
		//		log.Printf("%v\n", ok)
		//	}
		//	resp.(http.Flusher).Flush()
		//
		//	return nil
		//}))

		// 스트림을 열고 실행
		if err := <-strm.Open(); err != nil {
			errNotify <- struct{}{}
			resp.WriteHeader(http.StatusInternalServerError)
			log.Printf("Stream error :%s\n", err)
		}
	})

	log.Println("Server listening on: 4040")
	log.Fatal(http.ListenAndServe(":4040", nil))
}


func sendLogging(ch chan<- Log, errNotify <-chan struct{}, closeNotify <-chan bool) {
	defer close(ch)

	for i:=0; i<=5; i++ {
		select {
		case <-errNotify:
			log.Printf("err: stop goroutine")
			return
		case <-closeNotify:
			log.Printf("close: stop goroutine")
			return
		default:
			data := getLog()
			log.Printf("sending channel....%v\n", data)
			ch <- data
		}
	}
}


func getLog() Log{
	seed := rand.NewSource(time.Now().UnixNano())
	rn := rand.New(seed)
	size := len(Logs)

	return Logs[rn.Intn(size-1)]
}