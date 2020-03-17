package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// longRunTask runs after done is closed or timeout
func longRunTask(done <-chan interface{}, duration time.Duration) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case <-time.After(duration):
				ch <- rand.Int31n(10000000)
			}
		}
	}()
	return ch
}

func take(done <-chan interface{}, valueStream <-chan interface{}, n int) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for n > 0 {
			n--
			select {
			case <-done:
				return
			case ch <- <-valueStream:
			}
		}
	}()
	return ch
}

func fanIn(done <-chan interface{}, workers ...<-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})
	wg := &sync.WaitGroup{}

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for value := range c {
			select {
			case <-done:
				return
			case ch <- value:
			}
		}
	}
	wg.Add(len(workers))
	for _, w := range workers {
		go multiplex(w)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

func main() {
	done := make(chan interface{})
	defer close(done)

	// ========================== Sequential Run ================================
	// begin := time.Now()
	// sequentialExecution := longRunTask(done, time.Second)
	// for v := range take(done, sequentialExecution, 10) { // run longRunTask ten times sequentially, take 10s to finish.
	// 	fmt.Println(v)
	// }
	// fmt.Printf("done after: %v", time.Since(begin))

	// ============================= Parallel Run ================================
	N := runtime.NumCPU()
	workers := make([]<-chan interface{}, N)
	fmt.Printf("Run test on %v cores\n", N)
	for i := 0; i < N; i++ {
		workers[i] = longRunTask(done, time.Second*time.Duration(i+1))
	}
	parallelExecution := fanIn(done, workers...)
	begin := time.Now()
	for value := range take(done, parallelExecution, 10) { // run the longRunTask by N workers, and fanIn result to one channel
		fmt.Println(value)
	}
	fmt.Printf("done after: %v", time.Since(begin))
}
