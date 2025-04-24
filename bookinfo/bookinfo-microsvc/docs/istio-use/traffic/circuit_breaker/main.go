package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	concurrentRequest := 10
	url := "http://10.109.50.48/productpage?u=normal"
	for i := 0; i < concurrentRequest; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("request id %d failed, err %v\n", id, err)
				return
			}
			fmt.Printf("request id %d status code is %d\n", id, resp.StatusCode)
		}(i)
	}
	wg.Wait()
}
