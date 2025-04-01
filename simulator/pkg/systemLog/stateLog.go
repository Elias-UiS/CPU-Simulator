package systemLog

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/systemState"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type systemStateLog struct {
	PubSub *systemState.PubSub[systemState.State]
}

// data := systemState.State{Name: "Active", Age: 30}
// jsonData, err := json.Marshal(data)
// if err != nil {
//     panic(err)
// }
// fmt.Println(string(jsonData))

func NewSystemStateLog(pubSub *systemState.PubSub[systemState.State]) *systemStateLog {
	systemStateLog := &systemStateLog{
		PubSub: pubSub,
	}
	return systemStateLog
}

func (log systemStateLog) LogSystemState() {
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
		state = <-channel
		logger.Log.Println(state)
		writeToFile(file, state)
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
