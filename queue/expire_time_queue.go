package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func _worker(jobChan <- chan int)  {
	defer wg.Done()
	for job := range jobChan{
		fmt.Printf("执行任务 %d \n", job)
		time.Sleep(time.Second)
	}
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <- done:
		return true
	case <- time.After(timeout):
		return false
	}
}

func main()  {
	jobChan := make(chan int, 100)
	for i := 1; i <= 10; i++{
		jobChan <- i
	}
	wg.Add(1)
	close(jobChan)
	go _worker(jobChan)
	res := WaitTimeout(&wg, 5*time.Second)
	if res {
		fmt.Println("执行完成退出")
	} else {
		fmt.Println("执行超时退出")
	}

}
