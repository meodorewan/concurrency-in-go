// calculate the memory allocated for one goroutine
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func emptyGoroutine(wg *sync.WaitGroup, done <-chan interface{}) {
	wg.Done()
	<-done // block until main returned
}

func memConsumed() uint64 {
	runtime.GC()
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}

func main() {
	done := make(chan interface{})
	wg := &sync.WaitGroup{}
	N := 1e5
	wg.Add(int(N))

	begin := memConsumed()
	for i := N; i > 0; i-- {
		go emptyGoroutine(wg, done)
	}
	wg.Wait()
	end := memConsumed()
	close(done)
	fmt.Printf("Goroutine mem %.3fkb", float64(end - begin)/N/1000)
}
