package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

type result struct {
	voltage    float64
	temprature float64
}

func newResult(raw string) *result {
	r := result{}
	s := strings.Split(raw, " ")
	if len(s) < 2 {
		r.voltage = 666
		r.temprature = 666
	} else {
		r.voltage, _ = strconv.ParseFloat(s[0], 64)
		r.temprature, _ = strconv.ParseFloat(s[1], 64)
	}
	return &r
}

func connectToArduino() (serial.Port, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	for _, port := range ports {
		log.Printf("Found port: %v\n", port)
	}
	if len(ports) == 0 {
		return nil, errors.New("No serial ports found!")
	}
	//TODO: add config for choosing the port
	mode := &serial.Mode{
		BaudRate: 9600}
	port, err := serial.Open(ports[0], mode)
	if err != nil {
		return nil, errors.New("error while trying to connect")
	}
	log.Printf("Connected successfully to the arduino on port: %v\n", ports[0])
	return port, nil
}

func ProbeArduino(port serial.Port) (result, error) {
	raw := ""
	buff := make([]byte, 30)
	for {
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		if err != nil {
			log.Printf("failed to read data: %v", err)
			return result{}, errors.New("failed to read data")
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		raw = string(buff[:n])
		if strings.Contains(raw, "\n") {
			r, _ := regexp.Compile("[0-9]{1,2}[,.][0-9]{6} [0-9]{1,2}[,.][0-9]{6}")
			raw = r.FindString(raw)
			break
		}
	}
	return *newResult(raw), nil
}

func reconnectToArduino() serial.Port {
	timeout := 2 * time.Second
	lock := false
	port, err := connectToArduino()
	if err != nil {
		lock = true
	}
	for lock {
		port, err = connectToArduino()
		if err == nil {
			lock = false
		}
		time.Sleep(timeout)
		timeout *= 2
		log.Printf("trying to reconnect in %v", timeout)
	}
	return port

}
