package memory

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/settings"
	"fmt"
)

type FreeList struct {
	freeList           []bool // tracks free physical frames
	numberOfFreeFrames int    // free frames
}

// Sets the frame to be not free.
func (freeList *FreeList) AllocateFrame(framesNeeded int) ([]int, error) {
	if framesNeeded > freeList.numberOfFreeFrames {
		return nil, fmt.Errorf("Not enough free frames.")
	}
	list := []int{}
	logger.Log.Println("INFO: AllocateFrame() 1")
	gotten := 0
	for i := range freeList.freeList {
		if freeList.freeList[i] == true {
			list = append(list, i)
			gotten += 1
			freeList.freeList[i] = false
		}
		if framesNeeded == gotten {
			break
		}

	}
	logger.Log.Println("INFO: AllocateFrame() 2")
	logger.Log.Printf("INFO: AllocateFrame() %d ", gotten)
	freeList.numberOfFreeFrames -= gotten
	return list, nil
}

// Sets the frame to be free.
func (freeList *FreeList) DeallocateFrame(list []int) {
	for i := range list {
		index := list[i]
		freeList.freeList[index] = true
		freeList.numberOfFreeFrames += 1
	}
}

func NewFreeList() *FreeList {
	freeList := &FreeList{numberOfFreeFrames: settings.NumFrames}
	freeList.freeList = make([]bool, settings.NumFrames)
	for i := range freeList.freeList {
		freeList.freeList[i] = true
	}
	return freeList
}

var FreelistObject *FreeList = NewFreeList()
