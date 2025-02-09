package memory

import "fmt"

type FreeList struct {
	freeList           []bool // tracks free physical frames
	numberOfFreeFrames int    // free frames
}

// Sets the frame to be not free.
func (freeList *FreeList) allocateFrame(framesNeeded int) ([]int, error) {
	if framesNeeded > freeList.numberOfFreeFrames {
		return nil, fmt.Errorf("Not enough free frames.")
	}
	list := []int{}
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
	freeList.numberOfFreeFrames -= gotten
	return list, nil
}

// Sets the frame to be free.
func (freeList *FreeList) deallocateFrame(list []int) {
	for i := range list {
		index := list[i]
		freeList.freeList[index] = true
		freeList.numberOfFreeFrames += 1
	}
}

func newFreeList() *FreeList {
	freeList := &FreeList{numberOfFreeFrames: NumFrames}
	freeList.freeList = make([]bool, NumFrames)
	for i := range freeList.freeList {
		freeList.freeList[i] = true
	}
	return freeList
}

var FreelistObject *FreeList = newFreeList()
