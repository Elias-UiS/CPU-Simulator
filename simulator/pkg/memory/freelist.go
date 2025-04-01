package memory

import (
	"CPU-Simulator/simulator/pkg/logger"
	"CPU-Simulator/simulator/pkg/settings"
	"fmt"
)

type FreeList struct {
	FreeList           []bool // tracks free physical frames
	NumberOfFreeFrames int    // free frames
}

// Sets the frame to be not free.
func (freeList *FreeList) AllocateFrame(framesNeeded int) ([]int, error) {
	if framesNeeded > freeList.NumberOfFreeFrames {
		return nil, fmt.Errorf("Not enough free frames.")
	}
	list := []int{}
	logger.Log.Println("INFO: AllocateFrame() 1")
	gotten := 0
	for i := range freeList.FreeList {
		if freeList.FreeList[i] {
			list = append(list, i)
			gotten += 1
			freeList.FreeList[i] = false
		}
		if framesNeeded == gotten {
			break
		}

	}
	logger.Log.Printf("INFO: AllocateFrame(), Frames given: %d ", gotten)
	for _, value := range list {
		logger.Log.Printf("INFO: AllocateFrame(), Frame nr: %d ", value)
	}
	freeList.NumberOfFreeFrames -= gotten
	return list, nil
}

// Sets the frame to be free.
func (freeList *FreeList) DeallocateFrame(list []int) {
	for i := range list {
		index := list[i]
		freeList.FreeList[index] = true
		freeList.NumberOfFreeFrames += 1
	}
}

func NewFreeList() *FreeList {
	freeList := &FreeList{NumberOfFreeFrames: settings.NumFrames}
	freeList.FreeList = make([]bool, settings.NumFrames)
	for i := range freeList.FreeList {
		freeList.FreeList[i] = true
	}
	return freeList
}

var FreelistObject *FreeList = NewFreeList()
