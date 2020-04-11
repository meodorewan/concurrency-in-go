// stdout:
// new heartbeats 1.000969845
// new heartbeats 2.001589157
// new result 2.00167027
// new heartbeats 3.002476467
// new heartbeats 4.004912529
// new result 4.004950586
// new heartbeats 5.004823972
// new heartbeats 6.005248108
// new result 6.005293691
// new heartbeats 7.00490851
// new heartbeats 8.000137867
// new result 8.00018069
// new heartbeats 9.002155347
// cancelled
package main

import (
	"context"
	"fmt"
	"time"
)

func f(ctx context.Context, d time.Duration) (chan interface{}, chan interface{}) {
	resultCh := make(chan interface{})
	pulseCh := make(chan interface{})
	go func() {
		defer close(resultCh)
		defer close(pulseCh)

		t := time.NewTicker(d)
		w := time.NewTicker(d * 2) // simulating result coming

		sendPulse := func() {
			select {
			case pulseCh <- struct{}{}:
			default: // ignore if there is no listener
			}
		}

		sendResult := func() {
			select {
			case <-ctx.Done():
				return
			case <-t.C: // send heartbeat as well
				sendPulse()
			case resultCh <- struct{}{}:
				return
			}
		}

		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancelled")
				return
			case <-t.C:
				sendPulse()
			case <-w.C:
				sendResult()
			}
		}
	}()
	return resultCh, pulseCh
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	result, l := f(ctx, time.Second)
	start := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-result:
			fmt.Println("new result", time.Since(start).Seconds())
		case <-l:
			fmt.Println("new heartbeats", time.Since(start).Seconds())
		}
	}
}
