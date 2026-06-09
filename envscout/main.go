package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	base := "192.168.1."
	port := "1433"

	var wg sync.WaitGroup
	sem := make(chan struct{}, 120) // limita concorrência

	fmt.Println("Scanning for SQL Server (TCP 1433) on", base+"0/24 ...")

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", base, i)

		wg.Add(1)
		sem <- struct{}{}
		go func(host string) {
			defer wg.Done()
			defer func() { <-sem }()

			addr := net.JoinHostPort(host, port)
			conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
			if err == nil {
				_ = conn.Close()
				fmt.Println("✅ Found open 1433 on:", host)
			}
		}(ip)
	}

	wg.Wait()
	fmt.Println("Done.")
}