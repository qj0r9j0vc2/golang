package main

import (
	"fmt"
	"sync"
	"time"
)

var globalVariable int

func add(i int, mutex *sync.Mutex, wg *sync.WaitGroup) {
	mutex.Lock()
	globalVariable += i
	mutex.Unlock()

	time.Sleep(1 * time.Second)
	wg.Done()
}

func main() {
	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(100)

	var mutex sync.Mutex

	for j := 0; j < 100; j++ {
		go add(j, &mutex, &wg)
	}

	wg.Wait()
	delta := time.Now().Sub(startTime)
	fmt.Printf("Result is %d, done in %.3fs.\n", globalVariable, delta.Seconds())
}
