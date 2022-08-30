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

func main() {
	logFile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logFile)
	ArduinoPort := connectToArduino() // connect to arduino

	res1 := ProbeArduino(ArduinoPort)
	curRec = startRecording("abc", float32(res1.temprature), float32(res1.voltage))

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
			res := ProbeArduino(ArduinoPort)
			temperature := fmt.Sprintf("%v", res.temprature)
			volt := fmt.Sprintf("%v", res.voltage)
			if temperature != "666" { // ignore invalid data
				lastValidRes = res
				server.BroadcastToNamespace("/", "temp", temperature+","+volt) // send data to client
				if curRec.BatterySerial != "" {
					_ = curRec.CSVwriter.Write([]string{fmt.Sprint(time.Now().Unix()), temperature, volt})
					curRec.CSVwriter.Flush()
				}
			}
		}
	}()

	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
