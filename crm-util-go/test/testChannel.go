package test

import (
	"crm-util-go/common"
	"fmt"
	"time"
)

func sum10(s int, writeChan chan<- int) {
	sum := s + 10
	fmt.Println("sum: ", sum)

	writeChan <- sum // send sum to channel
}

func TestChannel() {

	fmt.Println("Start TestChannel")

	// Like maps and slices, channels must be created before use:
	ch := make(chan int)
	// closes a channel
	defer close(ch)

	go sum10(1, ch)
	go sum10(2, ch)

	fmt.Println("Start receive")
	// receive from ch
	s1, ok := <-ch
	s2 := <-ch
	fmt.Println("value from channel:", s1, ok)
	fmt.Println("value from channel:", s2)

	for i := 0; i < 10; i++ {
		go sum10(i, ch)
	}
	fmt.Println("Length", len(ch))

	for j := 0; j < 10; j++ {
		s := <-ch
		fmt.Println("value from channel:", s)
	}
}

type Data struct {
	Id string
}

type Response struct {
	TransID string
	Code    string
	Msg     string
	Id      string
}

func invokeAPI(data Data, writeRespChan chan<- Response) {
	var resp Response
	resp.TransID = common.NewUUID()
	resp.Code = "0"
	resp.Msg = "Success"
	resp.Id = data.Id

	writeRespChan <- resp
}

func TestInvokeAPIChannel() {

	fmt.Println("Start TestInvokeAPIChannel")
	defer fmt.Println("End TestInvokeAPIChannel")

	// Like maps and slices, channels must be created before use:
	respChan := make(chan Response)
	// closes a channel
	defer close(respChan)

	var dataList []Data
	for i := 0; i < 17; i++ {
		data := Data{}
		data.Id = common.IntToString(i)

		dataList = append(dataList, data)
	}

	fmt.Println("Size: ", len(dataList))

	countThread := 0
	maxThread := 10

	for _, dataRow := range dataList {
		go invokeAPI(dataRow, respChan)
		countThread++

		if countThread == maxThread {
			for j := 0; j < countThread; j++ {
				response := <-respChan
				fmt.Println("1. Output from channel:"+response.TransID, response.Code, response.Msg)
			}

			countThread = 0
		}
	}

	if countThread > 0 {
		for j := 0; j < countThread; j++ {
			response := <-respChan
			fmt.Println("2. Output from channel:"+response.TransID, response.Code, response.Msg)
		}
	}
}

func multiplyByTwo(readInputChan <-chan int, writeOutputChan chan<- int) {
	fmt.Println("Start goroutine multiplyByTwo")
	defer fmt.Println("End goroutine multiplyByTwo")
	for {
		num, isOpenChan := <-readInputChan
		if isOpenChan {
			result := num * 2
			writeOutputChan <- result
		} else {
			return
		}
	}
}

func TestGoRoutinesChannel() {
	fmt.Println("Start TestGoRoutinesChannel")
	defer fmt.Println("End TestGoRoutinesChannel")

	defer time.Sleep(2 * time.Second)

	inputChan := make(chan int)
	defer close(inputChan)

	outputChan := make(chan int)
	defer close(outputChan)

	// Create 3 `multiplyByTwo` goroutines.
	go multiplyByTwo(inputChan, outputChan)
	go multiplyByTwo(inputChan, outputChan)
	go multiplyByTwo(inputChan, outputChan)

	// Up till this point, none of the created goroutines actually do
	// anything, since they are all waiting for the `in` channel to
	// receive some data, we can send this in another goroutine
	go func() {
		inputChan <- 1
		inputChan <- 2
		inputChan <- 3
		inputChan <- 4
	}()

	// Now we wait for each result to come in
	fmt.Println(<-outputChan)
	fmt.Println(<-outputChan)
	fmt.Println(<-outputChan)
	fmt.Println(<-outputChan)
}

func fast(num int, out chan<- int) {
	result := num * 2
	time.Sleep(5 * time.Millisecond)
	out <- result
}

func slow(num int, out chan<- int) {
	result := num * 2
	time.Sleep(15 * time.Millisecond)
	out <- result
}

func TestSelectChannel() {
	fmt.Println("Start TestSelectChannel")
	defer fmt.Println("End TestSelectChannel")

	out1 := make(chan int)
	defer close(out1)

	out2 := make(chan int)
	defer close(out2)

	// we start both fast and slow in different goroutines with different channels
	go fast(2, out1)
	go slow(3, out2)

	// perform some action depending on which channel receives information first
	select {
	case res := <-out1:
		fmt.Println("fast finished, result:", res)
	case res := <-out2:
		fmt.Println("slow finished, result:", res)
	}
}
