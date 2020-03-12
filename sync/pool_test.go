// Result
// goos: darwin
// goarch: amd64
// pkg: github.com/meodorewan/concurrency-in-go/sync
// BenchmarkNetworkRequestWithPool-4   	    4700	   8083372 ns/op
// PASS
// ok  	github.com/meodorewan/concurrency-in-go/sync	55.531s
// => 8083372 ns/op vs 1005064514 ns/op
package sync

import (
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"testing"
)
//
// func init() {
// 	wg := startNetWorkDaemonWithPool()
// 	wg.Wait()
// }

func cachedConnection() *sync.Pool {
	p := &sync.Pool{
		New: connectToServer,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startNetWorkDaemonWithPool() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	p := cachedConnection()
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
			connection := p.Get()
			p.Put(connection)
			conn.Close()
		}
	}()
	return wg
}

func BenchmarkNetworkRequestWithPool(b *testing.B) {
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
