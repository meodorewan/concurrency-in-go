package main

import (
	"fmt"
	"time"
)

func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

func main() {
	sig := func(duration time.Duration) <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			defer close(ch)
			time.Sleep(duration)
		}()
		return ch
	}
	start := time.Now()
	<-Or(
		sig(time.Hour),
		sig(time.Minute),
		sig(time.Second*5),
		sig(time.Second*2),
		sig(time.Second*3),
	)
	fmt.Printf("done after: %v", time.Since(start))
}
