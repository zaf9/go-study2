package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 简单并发压测客户端，用于验证日志开销与稳定性（非生产级压测工具）
func main() {
	url := flag.String("url", "http://localhost:8080/", "目标 URL")
	concurrency := flag.Int("concurrency", 100, "并发数")
	total := flag.Int("requests", 1000, "总请求数")
	timeout := flag.Int("timeout", 5, "单个请求超时（秒）")
	flag.Parse()

	client := &http.Client{Timeout: time.Duration(*timeout) * time.Second}
	wg := sync.WaitGroup{}
	reqCh := make(chan int)
	var mu sync.Mutex
	var success, failure int
	var totalDur time.Duration

	// workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range reqCh {
				start := time.Now()
				resp, err := client.Get(*url)
				dur := time.Since(start)
				mu.Lock()
				totalDur += dur
				if err != nil {
					failure++
				} else {
					resp.Body.Close()
					if resp.StatusCode >= 200 && resp.StatusCode < 400 {
						success++
					} else {
						failure++
					}
				}
				mu.Unlock()
			}
		}()
	}

	// pump requests
	go func() {
		for i := 0; i < *total; i++ {
			reqCh <- i
		}
		close(reqCh)
	}()

	wg.Wait()

	fmt.Printf("Target: %s\n", *url)
	fmt.Printf("Concurrency: %d, Requests: %d\n", *concurrency, *total)
	fmt.Printf("Success: %d, Failure: %d\n", success, failure)
	if success > 0 {
		avg := time.Duration(int64(totalDur) / int64(success))
		fmt.Printf("Average latency (for successful requests): %v\n", avg)
	}
}
