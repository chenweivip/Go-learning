package main

import (
	"context"
	"fmt"
	"time"
)

var ctx context.Context
var cancel context.CancelFunc

func worker(jobChan <- chan int, ctx context.Context)  {
	for {
		select {
		case <- ctx.Done():
			return
		case job := <- jobChan:
			fmt.Printf("执行任务 %d \n", job)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	jobChan := make(chan int, 100)
	// //带有取消功能的 contex
	ctx, cancel = context.WithCancel(context.Background())
	for i := 1; i <= 10; i++{
		jobChan <- i
	}
	close(jobChan)
	go worker(jobChan, ctx)
	time.Sleep(2 * time.Second)
	cancel()
}