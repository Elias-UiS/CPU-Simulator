package systemLog

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/systemState"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type SystemStateLog struct {
	PubSub    *systemState.PubSub[systemState.State]
	StopChan  chan bool
	StopState bool
}

// data := systemState.State{Name: "Active", Age: 30}
// jsonData, err := json.Marshal(data)
// if err != nil {
//     panic(err)
// }
// fmt.Println(string(jsonData))

func NewSystemStateLog(pubSub *systemState.PubSub[systemState.State]) *SystemStateLog {
	stopChan := make(chan bool)
	systemStateLog := &SystemStateLog{
		PubSub:    pubSub,
		StopChan:  stopChan,
		StopState: false,
	}
	return systemStateLog
}

func (log *SystemStateLog) LogSystemState() {
	for {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("logStateFiles/state_%s.json", timestamp)
		file, err := os.Create(filename)
		if err != nil {
			return
		}
		defer file.Close()

		channel := log.PubSub.Subscribe()
		var state systemState.State
		for {
			select {
			case <-log.StopChan:
				log.StopState = true
				logger.Log.Println("INFO: LogSystemState() - Stop signal received, exiting.")
				return
			case state = <-channel:
				logger.Log.Println(state)
				writeToFile(file, state)
			}

		}
	}
}

func writeToFile(file *os.File, state systemState.State) error {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(state)
	if err != nil {
		return err
	}

	// Write the JSON data to the file, followed by a newline for better readability
	_, err = file.Write(append(jsonData, '\n'))
	if err != nil {
		return err
	}

	return nil
}
