package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type recording struct {
	BatterySerial  string
	ChargeDuration time.Duration
	StartTime      time.Time
	EndTime        time.Time
	StartVolt      float32
	EndVolt        float32
	StartTemp      float32
	EndTemp        float32
	LogPath        string
	CSVwriter      *csv.Writer
}

func GetFilenameDate(r recording) string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	// Use layout string for time format.
	const layout = "01-02-2006"
	// Place now in the string.
	fname := r.BatterySerial + "-" + r.StartTime.Format(layout) + ".csv"
	return filepath.Join(path, "data", fname)
}

func startRecording(batterySerial string, startTemp float32, startVolt float32) *recording {
	r := recording{}
	r.BatterySerial = batterySerial
	AddBatterySerial(batterySerial)
	r.StartTemp = startTemp
	r.StartVolt = startVolt
	r.StartTime = time.Now()
	r.LogPath = GetFilenameDate(r)
	csvfile, err := os.OpenFile(r.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	r.CSVwriter = csv.NewWriter(csvfile)
	log.Printf("starting to record: %v (%v volts) to: %v", r.BatterySerial, r.StartVolt, r.LogPath)
	fmt.Printf("starting to record: %v (%v volts) to: %v \n", r.BatterySerial, r.StartVolt, r.LogPath)
	return &r
}

func finishRecording(r recording, endTemp float32, endVolt float32) {

	r.EndTemp = endTemp
	r.EndVolt = endVolt
	r.EndTime = time.Now()
	r.ChargeDuration = r.EndTime.Sub(r.StartTime)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("./recordings.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("failed to log recording:", err)
	}
	if _, err := f.Write([]byte(formatRecording(r))); err != nil {
		log.Println("failed to log recording:", err)
	}
	if err := f.Close(); err != nil {
		log.Println("failed to log recording:", err)
	}
	log.Printf("finished to record: %v after: %v", r.BatterySerial, r.ChargeDuration)
	fmt.Printf("finished to record: %v after: %v\n", r.BatterySerial, r.ChargeDuration)
}

func formatRecording(r recording) string {
	st := ""
	const layout = "2006-01-02 15:04:05"
	// start time
	st += r.StartTime.Format(layout) + ","
	// battery serial
	st += r.BatterySerial + ","
	// charging duration
	st += fmt.Sprintf("%s,", r.ChargeDuration.Round(time.Second))
	// values at start and end
	st += fmt.Sprintf("%v,", r.StartVolt)
	st += fmt.Sprintf("%v,", r.EndVolt)
	st += fmt.Sprintf("%v,", r.StartTemp)
	st += fmt.Sprintf("%v,", r.EndTemp)
	// end time
	st += r.EndTime.Format(layout) + ","
	// logfile
	st += r.LogPath
	// newline
	st += "\n"
	return st
}
