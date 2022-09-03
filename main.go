package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

var curRec *recording       // current recording
var lastValidRes = result{} // last valid sensor readings

func newRecordingRequest(w http.ResponseWriter, r *http.Request) {
	// handels request for a new recording
	if r.Method == "POST" {
		batterySerial := r.FormValue("serial") // get serial from the request
		curRec = startRecording(batterySerial, float32(lastValidRes.temprature), float32(lastValidRes.voltage))
	}
}

func endRecordingRequest(w http.ResponseWriter, r *http.Request) {
	// handels request for end of recording
	finishRecording(*curRec, float32(lastValidRes.temprature), float32(lastValidRes.voltage))
	curRec = nil
}

func getStateRequest(w http.ResponseWriter, r *http.Request) {
	// in case of browser refresh the client might lose sync, this function keep it in sync with the server
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		data := make(map[string]interface{})
		if curRec == nil {
			data = map[string]interface{}{
				"isRecording": false,
				"serials":     getBatterySerials(),
			}
		} else {
			data = map[string]interface{}{
				"isRecording":        true,
				"batterySerial":      curRec.BatterySerial,
				"startRecordingTime": curRec.StartTime.Unix(),
				"serials":            getBatterySerials(),
			}
		}

		jsonData, err := json.Marshal(data)
		b, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%s\n", b)
		// set header to 'application/json'
		w.Header().Set("Content-Type", "application/json")
		// write the response
		w.Write(b)
	}
}

func getBatterySerials() []string {
	content, err := ioutil.ReadFile("battery-serials.txt")
	if err != nil {
		log.Println(err)
		return []string{}
	}
	slc := strings.Split(string(content), "\n")
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}
	return slc
}

func AddBatterySerial(curSerial string) {
	serials := getBatterySerials()
	for _, v := range serials {
		if v == curSerial {
			return
		}
	}
	f, err := os.OpenFile("./battery-serials.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	if _, err := f.Write([]byte(fmt.Sprintf("%+v\n", curSerial))); err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
}

func main() {
	logFile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644) // log to file
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
			time.Sleep(time.Second / 3)           // sample arduino in 3Hz
			res, err := ProbeArduino(ArduinoPort) // try to probe arduino
			if err != nil {
				log.Println(err)
				ArduinoPort = reconnectToArduino() // in case of error try to reconnect to arduino
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
	http.HandleFunc("/getstate", getStateRequest)
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
