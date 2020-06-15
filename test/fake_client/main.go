package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//dial tcp 127.0.0.1:8080: connectex: No connection could be made because the target machine actively refused it.

func main() {
	//bodyStr := `["http://localhost:8081/1", "http://localhost:8081/2", "http://localhost:8081/3", "http://localhost:8081/4", "http://localhost:8081/5", "http://localhost:8081/6", "http://localhost:8081/7", "http://localhost:8081/8", "http://localhost:8081/9", "http://localhost:8081/10", "http://localhost:8081/11", "http://localhost:8081/12", "http://localhost:8081/13", "http://localhost:8081/14", "http://localhost:8081/15", "http://localhost:8081/16", "http://localhost:8081/17", "http://localhost:8081/18", "http://localhost:8081/19", "http://localhost:8081/20", "http://localhost:8081/21"]`
	//bodyStr := `["http://localhost:8081/1", "http://localhost:8081/2", "http://localhost:8081/3", "http://localhost:8081/4", "http://localhost:8081/5", "http://localhost:8081/6", "http://localhost:8081/7", "http://localhost:8081/8", "http://localhost:8081/9", "http://localhost:8081/10", "http://localhost:8081/11", "http://localhost:8081/12", "http://localhost:8081/13", "http://localhost:8081/14", "http://localhost:8081/15", "http://localhost:8081/16", "http://localhost:8081/17", "http://localhost:8081/18", "http://localhost:8081/19", "http://localhost:8081/20"]`
	//bodyStr := `["http://localhost:8081/1", "http://localhost:8081/2", "http://localhost:8081/3"]`
	bodyStr := `["https://ya.ru"]`

	wg := &sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < 1; i++ {
		go func(wg *sync.WaitGroup, i int) {
			req, err := http.NewRequest("POST", "http://localhost:8080/", bytes.NewReader([]byte(bodyStr)))
			if err != nil {
				fmt.Println(err)
			}
			ctx, cancel := context.WithTimeout(req.Context(), 600*time.Millisecond)
			//ctx, cancel := context.WithCancel(req.Context())
			defer cancel()

			req = req.WithContext(ctx)
			client := &http.Client{}
			fmt.Println(time.Now())
			resp, err := client.Do(req)
			fmt.Println(time.Now())

			if err != nil {
				fmt.Println(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(i)
			fmt.Println("Status:", resp.Status)
			fmt.Println("Body:", string(body))
			wg.Done()
		}(wg, i)
	}
	wg.Wait()
}
