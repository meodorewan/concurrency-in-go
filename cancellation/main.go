package main

import (
	"fmt"
	"sync"
	"time"
)

func doReadWork(done, data <-chan interface{}) <-chan interface{} {
	terminated := make(chan interface{})
	go func() {
		defer close(terminated)
		defer fmt.Println("goroutine closed")
		for {
			select {
			case <-done:
				fmt.Println("be cancelled")
				return
			case <-data:
			}
		}
	}()
	return terminated
}

func testReadChan() {
	done := make(chan interface{})
	terminated := doReadWork(done, nil)
	go func() {
		time.Sleep(time.Second * 2)
		close(done)
	}()
	<-terminated
}

func main() {
	// testReadChan()
	testWriteChan()
}

func testWriteChan() {
	done := make(chan interface{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	data := doWriteWork(done, wg)
	for i := 0; i < 10; i++ {
		fmt.Println(<-data)
	}
	close(done)
	wg.Wait()
	fmt.Println("done")
}

func doWriteWork(done chan interface{}, wg *sync.WaitGroup) <-chan int {
	result := make(chan int)
	go func() {
		defer wg.Done()
		defer close(result)
		cnt := 0
		for {
			cnt++
			select {
			case <-done:
				fmt.Println("be cancelled")
				return
			case result <- cnt:
			}
		}
	}()
	return result
}
