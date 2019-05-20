package main

import (
	"fmt"
	"runtime"
	"time"
	"reflect"
)

var (
	MAX_WORKER = 10
)

type Payload struct {
	Num int
}

type Job struct {
	payload Payload
}

var JobQueue chan Job

// 执行任务的工作单元
type Worker struct {
	WorkerPool 	chan chan Job  // 工作者池--每个元素是一个工作者的私有任务channal
	JobChannel 	chan Job	 // 每个工作者单元包含一个任务管道 用于获取任务
	quit 		chan bool   // 退出信号
	no 			int  // 编号
}

// 创建一个新的工作单元
func NewWorker(workerPool chan chan Job, no int) *Worker {
	fmt.Println("创建一个新工作者单元")
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit: make(chan bool),
		no: no,
	}
}

func (w *Worker) Start()  {
	go func() {
		w.WorkerPool <- w.JobChannel
		fmt.Println("w.WorkerPool <- w.JobChannel", w)
		select {
		case job := <- w.JobChannel:
			fmt.Println("job := <-w.JobChannel")
			// 收到任务
			fmt.Println(job)
			time.Sleep(100 * time.Second)
		case <- w.quit:
			return
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	// 工作池
	WokerPool 	chan chan Job
	// 工作者数量
	MaxWorkers 	int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WokerPool: make(chan chan Job, maxWorkers),
		MaxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) dispatch()  {
	for  {
		select {
		case job := <- JobQueue:
			fmt.Println("job := <-JobQueue:")
			go func(job Job) {
				fmt.Println("等待空闲worker (任务多的时候会阻塞这里")
				// 等待空闲worker (任务多的时候会阻塞这里)
				jobChannel := <- d.WokerPool
				fmt.Println("jobChannel := <-d.WorkerPool", reflect.TypeOf(jobChannel))
				// 将任务放到上述woker的私有任务channal中
				jobChannel <- job
				fmt.Println("jobChannel <- job")
			}(job)
		}
	}
}

func (d *Dispatcher) Run()  {
	for i := 1; i < d.MaxWorkers+1; i++  {
		worker := NewWorker(d.WokerPool, i)
		worker.Start()
	}
	go d.dispatch()
}

func addQueue()  {
	for i := 0; i < 20; i++{
		// 新建一个任务
		payLoad := Payload{1}
		work := Job{payLoad}

		// 任务放到任务队列channel
		JobQueue <- work

		fmt.Println("JobQueue <- work", i)
		fmt.Println("当前协程数:", runtime.NumGoroutine())
		time.Sleep(100 * time.Millisecond)
	}
}

func main()  {
	JobQueue = make(chan Job, 10)
	dispatcher := NewDispatcher(MAX_WORKER)
	dispatcher.Run()
	time.Sleep(time.Second)
	go addQueue()
	time.Sleep(1000 * time.Second)
}
