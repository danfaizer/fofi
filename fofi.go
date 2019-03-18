package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// concurrentGoRoutines defines the # of concurrent Producer go routines
	// running at the same time.
	concurrentGoRoutines = 5
)

// Producer simulates a slow and resource consuming job.
func Producer(input string, out chan int, lock chan bool, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		<-lock
	}()
	lock <- true
	fmt.Printf("producer processing %s\n", input)

	time.Sleep(2 * time.Second)

	i, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("error: producer couldn't process %s\n", input)
		return
	}
	out <- i
	fmt.Printf("producer wrote %s\n", input)
}

// Consumer simulates a fast and low resource consuming job.
func Consumer(in chan int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for i := range in {
		fmt.Printf("consumer processing %d\n", i)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("consumer processed %d\n", i)
	}
	fmt.Println("consumer done")
}

func main() {
	sampleInputData := []string{"1", "2", "3", "4", "5"}

	var consumerWg sync.WaitGroup
	var producerWg sync.WaitGroup
	consumerCh := make(chan int, concurrentGoRoutines)
	producerCh := make(chan bool, concurrentGoRoutines)

	consumerWg.Add(1)
	go Consumer(consumerCh, &consumerWg)

	fmt.Println("starting fofi processing")
	for _, str := range sampleInputData {
		producerWg.Add(1)
		go Producer(str, consumerCh, producerCh, &producerWg)
	}

	producerWg.Wait()
	close(consumerCh)
	close(producerCh)
	consumerWg.Wait()
	fmt.Println("fofi processing finished")
}
