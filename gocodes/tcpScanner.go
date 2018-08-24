package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Ip: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	target := os.Args[1]

	/*sync.WaitGroup provides a goroutine synchronization mechanism in Golang, and is used for waiting for a collection of goroutines to finish*/
	wg := sync.WaitGroup{}

	c := func(ports int) {
		conn, err := net.DialTimeout("tcp", target+":"+strconv.Itoa(ports), time.Duration(1)*time.Second)
		if err == nil {
			fmt.Fprintf(conn, "HELLO\r\n")
			buf := make([]byte, 0, 4096) // size of buffer
			tmp := make([]byte, 256)     // small buffer tmp
			for {
				n, err := conn.Read(tmp)
				if err != nil {
					if err != io.EOF {
						fmt.Println("read error", err)
					}
					break
				}
				buf = append(buf, tmp[:n]...)
			}
			conn.Close()
			fmt.Println(ports, " open")

		}
		wg.Done()
	}
	wg.Add(65534)
	for i := 0; i < 65534; i++ {
		go c(i)
	}
	wg.Wait()
	os.Exit(0)
}
