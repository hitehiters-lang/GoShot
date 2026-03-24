package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.design/x/clipboard"
)

const version = "1.0"

type Configuration struct {
	SaveDir       string `json:"saveDir"`
	CheckInterval int    `json:"checkInterval"`
}

func main() {

	conf := Configuration{
		SaveDir:       "..\\screenshots",
		CheckInterval: 1,
	}

	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	logFile := filepath.Join("goshot.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("open log error: ", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime)

	log.Println("===> GoShot", version, "launched at", execDir)

	configPath := filepath.Join(execDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "config.json"
		log.Println("got config")
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalln("config open error:", err)
	}
	defer file.Close()
	log.Println("opened config")

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatalln("config parsing error:", err)
	}
	log.Println("saveDir path readed from config:", conf.SaveDir)

	if err = os.MkdirAll(conf.SaveDir, os.ModePerm); err != nil {
		log.Fatalln("mkdir error:", err)
	}

	var lastSize int
	ticker := time.NewTicker(time.Second * time.Duration(conf.CheckInterval))
	defer ticker.Stop()

	if err := clipboard.Init(); err != nil {
		log.Fatalln("clipboard init error:", err)
	}
	log.Println("clipboard inited")

	for {
		<-ticker.C

		data := clipboard.Read(clipboard.FmtImage)
		if data == nil {
			continue
		}

		if len(data) == lastSize {
			continue
		}
		lastSize = len(data)

		filename := filepath.Join(conf.SaveDir, fmt.Sprintf("goshot_%s.png", time.Now().Format("20060102_150405")))

		if err := os.WriteFile(filename, data, 0644); err != nil {
			log.Println("write file error:", err)
			continue
		}

		log.Printf("saved: %s (%.2f KB)\n", filepath.Base(filename), float64(len(data))/1024)
	}
}

// building with
// go build -ldflags="-H windowsgui" -o GoShot.exe
