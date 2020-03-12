// Result
// goos: darwin
// goarch: amd64
// pkg: github.com/meodorewan/concurrency-in-go/sync
// BenchmarkNetworkRequest-4           	      10	1005064514 ns/op
// BenchmarkNetworkRequestWithPool-4   	      10	1005235238 ns/op
// PASS
// ok  	github.com/meodorewan/concurrency-in-go/sync	23.454s
package sync

import (
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"
)

func init() {
	wg := startNetWorkDaemon()
	wg.Wait()
}

func connectToServer() interface{} {
	time.Sleep(time.Second)
	return struct {}{}
}

func startNetWorkDaemon() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8000")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			connectToServer()
			conn.Close()
		}
	}()
	return wg
}

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			b.Fatalf("can not dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("can not read: %v", err)
		}
		conn.Close()
	}
}
