package scanner

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func ScanOpenPort(base string, port int, timeoutMs int, concurrency int) []string {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []string
		sem     = make(chan struct{}, concurrency)
	)

	timeout := time.Duration(timeoutMs) * time.Millisecond
	portStr := fmt.Sprintf("%d", port)

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", base, i)

		wg.Add(1)
		sem <- struct{}{}
		go func(host string) {
			defer wg.Done()
			defer func() { <-sem }()

			addr := net.JoinHostPort(host, portStr)
			conn, err := net.DialTimeout("tcp", addr, timeout)
			if err == nil {
				_ = conn.Close()
				mu.Lock()
				results = append(results, host)
				mu.Unlock()
			}
		}(ip)
	}

	wg.Wait()
	sort.Strings(results)
	return results
}
