package test

import (
	"crm-util-go/common"
	"github.com/go-co-op/gocron"
	"fmt"
	"time"
)

func TestStartScheduler() {
	thaiLocation, _ := time.LoadLocation("Asia/Bangkok")

	s := gocron.NewScheduler(thaiLocation)

	// SingletonMode() -> a long running job will not be rescheduled until the current run is completed
	job1, _ := s.Every(2).Second().Do(Job1)
	job1.SingletonMode()

	job2, _ := s.Every(4).Second().Do(Job2)
	job2.SingletonMode()


	// you can start running the scheduler in two different ways:
	// starts the scheduler asynchronously
	// s.StartAsync()
	// fmt.Printf("Job: %v, Error: %v", job, err)

	// starts the scheduler and blocks current execution path
	s.StartBlocking()
}

func Job1() {
	transID := common.NewUUID()

	// Logic Program
	currDateTime := time.Now().Local()
	yyyymmddhhmmss := currDateTime.Format("2006-01-02T15:04:05.000+0700")
	fmt.Printf("TransID: %s, Job1: %s \n", transID, yyyymmddhhmmss)
}

func Job2() {
	transID := common.NewUUID()

	// Logic Program
	currDateTime := time.Now().Local()
	yyyymmddhhmmss := currDateTime.Format("2006-01-02T15:04:05.000+0700")
	fmt.Printf("TransID: %s, Job2: %s \n", transID, yyyymmddhhmmss)
}