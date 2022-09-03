package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

var curRec *recording
var lastValidRes = result{}

func sendRandom() float64 {
	return (rand.Float64() * 5) + 5
}

func newRecordingRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		batterySerial := r.FormValue("serial")
		curRec = startRecording(batterySerial, float32(lastValidRes.temprature), float32(lastValidRes.voltage))
	}
}

func endRecordingRequest(w http.ResponseWriter, r *http.Request) {
	finishRecording(*curRec, float32(lastValidRes.temprature), float32(lastValidRes.voltage))
	curRec = nil
}

func main() {
	logFile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	} else {
		log.SetOutput(logFile)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ArduinoPort := reconnectToArduino() // connect to arduino

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()

	go func() {
		for {
			time.Sleep(time.Second / 3) // sample arduino in 4Hz
			res, err := ProbeArduino(ArduinoPort)
			if err != nil {
				log.Println(err)
				ArduinoPort = reconnectToArduino()
				res = result{}
			}
			temperature := fmt.Sprintf("%v", res.temprature)
			volt := fmt.Sprintf("%v", res.voltage)
			if temperature != "666" { // ignore invalid data
				lastValidRes = res
				server.BroadcastToNamespace("/", "temp", temperature+","+volt) // send data to client
				if curRec != nil {
					_ = curRec.CSVwriter.Write([]string{fmt.Sprint(time.Now().UnixMilli()), temperature, volt})
					curRec.CSVwriter.Flush()
				}
			}
		}
	}()

	defer server.Close()

	http.HandleFunc("/newrec", newRecordingRequest)
	http.HandleFunc("/endrec", endRecordingRequest)
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
