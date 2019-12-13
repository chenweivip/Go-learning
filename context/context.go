package main

import (
	"errors"
	"sync"
)
import "context"

func Rpc(ctx context.Context, url string) error  {
	result := make(chan int)
	err := make(chan error)

	go func() {
		isSucess := true
		if isSucess {
			result <- 1
		} else {
			err <- errors.New("some error happen")
		}
	}()

	select {
	case <- ctx.Done():
		return ctx.Err()
	case e := <- err:
		return e
	case <- result:
		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	err := Rpc(ctx, "http://rpc_1_url")
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}

	// RPC2调用
	wg.Add(1)
	go func(){
		defer wg.Done()
		err := Rpc(ctx, "http://rpc_2_url")
		if err != nil {
			cancel()
		}
	}()

	// RPC3调用
	wg.Add(1)
	go func(){
		defer wg.Done()
		err := Rpc(ctx, "http://rpc_3_url")
		if err != nil {
			cancel()
		}
	}()

	// RPC4调用
	wg.Add(1)
	go func(){
		defer wg.Done()
		err := Rpc(ctx, "http://rpc_4_url")
		if err != nil {
			cancel()
		}
	}()

	wg.Wait()
}

