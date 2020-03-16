package main

import (
	"fmt"
)

func add(done, input <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case num := <-input:
				ch <- num + 1
			}
		}
	}()
	return ch
}

func square(done, input <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case num := <-input:
				ch <- num * num
			}
		}
	}()
	return ch
}

func createPipeline(done chan int, data []int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		defer close(done)
		for _, num := range data {
			ch <- num
		}
	}()
	return ch
}

func main() {
	done := make(chan int)
	ch := createPipeline(done, []int{1, 2, 3})
	added := add(done, ch)
	squared := square(done, added)
	for num := range squared {
		fmt.Println(num)
	}
}
