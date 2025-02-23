// pkg/logger/logger.go
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var Log *log.Logger

func Init() {
	currentTime := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("logFiles/log_%s.txt", currentTime)
	if err := os.MkdirAll("logFiles", 0755); err != nil {
		log.Fatal("Could not create logs directory:", err)
	}
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Create MultiWriter for console and file logging
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	Log = log.New(multiWriter, "", log.LstdFlags)
}
