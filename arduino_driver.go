package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

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

func connectToArduino() serial.Port {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	for _, port := range ports {
		log.Printf("Found port: %v\n", port)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	//TODO: add config for choosing the port
	mode := &serial.Mode{
		BaudRate: 9600}
	port, err := serial.Open(ports[0], mode)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected successfully to the arduino on port: %v\n", ports[0])
	return port
}

func ProbeArduino(port serial.Port) result {
	raw := ""
	buff := make([]byte, 30)
	for {
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
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
	return *newResult(raw)
}
