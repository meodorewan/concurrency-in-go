package main

import (
	"fmt"
	"time"
)

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

func tee(done, valueStream <-chan interface{}) (_, _ chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for {
			select {
			case <-done:
				return
			case v, ok := <-valueStream:
				if !ok {
					return
				}
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
						return
					case out1 <- v:
						out1 = nil // to send v to out2 in the next iterator
					case out2 <- v:
						out2 = nil // to send v to out1 in the next iterator
					}
				}
			}
		}
	}()
	return out1, out2
}

func main() {
	generator := func(done <-chan interface{}) <-chan interface{} {
		ch := make(chan interface{})
		cnt := 0
		go func() {
			defer close(ch)
			for {
				cnt += 1
				select {
				case <-done:
					return
				case ch <- cnt:
				}
			}
		}()
		return ch
	}

	done := make(chan interface{})
	go func() {
		select {
		case <-time.After(time.Second * 10):
			close(done)
		}
	}()

	valueStream := take(done, generator(done), 10)
	out1, out2 := tee(done, valueStream)
	for v := range out1 {
		fmt.Printf("value from out1: %v\n", v)
		v2 := <-out2
		fmt.Printf("value from out2: %v\n", v2)
		fmt.Println("============")
	}
}
