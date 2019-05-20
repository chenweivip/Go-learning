package main

import (
	"fmt"
	"time"
)

func worker(jobChan <-chan int, cancelChan <-chan struct{}) {
	for {
		select {
		case <-cancelChan:
			return
		case job := <-jobChan:
			fmt.Printf("执行任务 %d \n", job)
			time.Sleep(1 * time.Second)
		}
	}
}

func main()  {
	jobChan := make(chan int, 100)
	cancelChan := make(chan struct{})

	for i := 1; i <= 10; i++ {
		jobChan <- i
	}

	close(jobChan)
	go worker(jobChan, cancelChan)
	time.Sleep(2 * time.Second)
	close(cancelChan)

}
