package test

import "fmt"

func myPanic(s string) {
	panic(s)      // throws panic
}

func myRecover() {
	e := recover()

	if e != nil {
		fmt.Println("Recovered from panic")
	}
}

func TestPanic(msg string) {
	defer myRecover()
	myPanic(msg)
}
