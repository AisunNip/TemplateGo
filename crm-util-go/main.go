package main

import (
	"crm-util-go/sseapi"
	"fmt"
	"os"
)

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func TestEnvironmentVariables() {
	// Environment Variables
	os.Setenv("FirstName", "Paravit")

	fmt.Println("FirstName:", os.Getenv("FirstName"))
	fmt.Println("JAVA_HOME:", os.Getenv("JAVA_HOME"))
	fmt.Println("Path:", os.Getenv("Path"))
}

func main() {
	/*
		go test crm-util-go/benchmark -bench=.
		go test crm-util-go/benchmark -bench=BenchmarkIOReadAll
		go test crm-util-go/benchmark -bench=BenchmarkIOCopy
		go test crm-util-go/benchmark -bench=BenchmarkBMHasStringValue*

		go test -v crm-util-go/array
		go test -v crm-util-go/cryptography
		go test -v crm-util-go/logging
		go test -v crm-util-go/errorcode
		go test -v crm-util-go/httpclient
		go test -v crm-util-go/pointer
		go test -v crm-util-go/validate
		go test -v crm-util-go/crmencoding
		go test -v crm-util-go/db
		go test -v crm-util-go/db -run "TestDeleteEmployeeMongo"
		go test -v crm-util-go/qrcode
		go test -v crm-util-go/timeUtil
	*/

	/*
		for i := 1; i < 10; i++ {
			key, _ := common.GenerateUniqueKey("CRMP")
			fmt.Println(key)
		}
	*/

	//httpclient.StartHttpProxy(":80")

	// HTTP SSE API
	streamIDList := []string{"m1", "m2"}
	sseapi.StartSSEApi(":8080", streamIDList)

	// test.StartEchoHTTP2Server()

	//google.StartHttp()

	// Go Routines and Channel
	//test.TestGoRoutinesChannel()
	//test.TestInvokeAPIChannel()
	//test.TestSelectChannel()

	// test.TestCassandraProduction()

	// TestEnvironmentVariables()

	// test.TestPanic("test panic!!")

	// test.TestSortStruct()

	// test.TestStartScheduler()

	// test.TestDateTime()

	// test.TestReflect(20.001)

	// test.TestFile()

	// test.TestSubslicingArray()

	// Java use "AES/ECB/PKCS5Padding"
	// test.TestAESGcm()
	// test.TestAESEcb()
	// test.TestHashMessage()

	// test.TestSendEmail()

	// test.TestSignal()

	//x := time.Now().Unix()
	//fmt.Println(x)
	//
	//y := time.Unix(x, 0)
	//fmt.Println(y)

	// test.TestChannel()

	//var x string
	//x = "hello 5555"
	//
	//rsMap := make(map[string]interface{})
	//rsMap["srOwner"] = &x
	//
	//y := rsMap["srOwner"].(*string)
	//fmt.Println(*y)

	/*
		fmt.Println(">>> 1. Sleep: ")
		time.Sleep(time.Duration(10) * time.Second)

		for i := 0; i < 20; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				transID, _ := common.NewUUID()
				campTransRespBean := test.GetCampTransListDAO(transID, campReqBean)
				fmt.Println(">>> Code: ", campTransRespBean.Code, "Msg:", campTransRespBean.Msg)
			}()
		}

		wg.Wait()


		fmt.Println(">>> 2. Sleep: ")
		time.Sleep(time.Duration(10) * time.Second)

		for i := 0; i < 20; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				transID, _ := common.NewUUID()
				campTransRespBean := test.GetCampTransListDAO(transID, campReqBean)
				fmt.Println(">>> Code: ", campTransRespBean.Code, "Msg:", campTransRespBean.Msg)
			}()
		}

		wg.Wait()

		test.CloseDB("100")
	*/

	/*
		currentDateTime := time.Now()
		fmt.Println(currentDateTime)

		switch {
			case currentDateTime.Hour() < 12:
				fmt.Println("Good morning!")
			case currentDateTime.Hour() < 17:
				fmt.Println("Good afternoon.")
			default:
				fmt.Println("Good evening.")
		}

		s := []int{2, 3, 5, 7, 11, 13}
		printSlice(s)

		// Slice the slice to give it zero length.
		s = s[:0]
		printSlice(s)

		// Extend its length.
		s = s[:4]
		printSlice(s)

		// Drop its first two values.
		s = s[2:]
		printSlice(s)

		var book1 Books
		book1.Title = "abc"
		book1.Author = "Paravit"

		var book2 Books
		book2.Title = "xyz"

		// new(type) --> return pointer
		var book3 = new(Books)
		book3.Title = "3333"

		var book4 = new(Books)
		book4.Title = "44444"

		var book5 *Books
		book5 = book4
		book5.Title = "5555"

		fmt.Println(book1)
		fmt.Println(book2)
		fmt.Println(book3)
		fmt.Println(book4)
		fmt.Println(book5)

		jsonString, _ := crmencoding.StructToJsonIndent(book1)
		fmt.Println(jsonString)

		// Unmarshal
		// jsonString := `{"title":"abc","author":"","subject":"","bookID":0}`
		var bookx Books
		err := json.Unmarshal([]byte(jsonString), &bookx)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(bookx)
			fmt.Println(bookx.Title)
		}
	*/
}
