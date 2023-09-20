package test

import (
	"fmt"
	"sync"
)

func logic(wg *sync.WaitGroup, data []int) {
	defer wg.Done()

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(data)
}

func TestSubslicingArray() {
	dataList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21,22,23,24,25,26}
	// dataList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,20}

	totalRecord := len(dataList)

	limitRecord, noOfThread := GetNoOfThread(totalRecord)

	fmt.Println("totalRecord", totalRecord)
	fmt.Println("limitRecord", limitRecord)
	fmt.Println("noOfThread", noOfThread)

	var wg sync.WaitGroup

	for lowerBoundIndex := 0; lowerBoundIndex < totalRecord; lowerBoundIndex += limitRecord {
		subslicingData := dataList[lowerBoundIndex:getUpperBoundIndex(lowerBoundIndex + limitRecord, totalRecord)]

		wg.Add(1)
		go logic(&wg, subslicingData)
	}

	wg.Wait()

	fmt.Println("Finish TestSubslicingArray")
}

func getUpperBoundIndex(a int, b int) int {
	if a <= b {
		return a
	}

	return b
}

func GetNoOfThread(totalRecord int) (limitRecord int, noOfThread int) {
	maxThread := 12

	if totalRecord <= maxThread {
		noOfThread = 1
	} else {
		noOfThread = totalRecord / maxThread

		if noOfThread > maxThread {
			noOfThread = maxThread
		}

		fraction := totalRecord % maxThread

		if fraction > 0 {
			noOfThread++
		}
	}

	limitRecord = totalRecord / noOfThread

	noOfThread = 0

	for lowerBoundIndex := 0; lowerBoundIndex < totalRecord; lowerBoundIndex += limitRecord {
		noOfThread++
	}

	return limitRecord, noOfThread
}