package main

import (
	"crypto/sha1"
	"fmt"
	"time"
	"sync"
	"runtime"
	"math/rand"
)

type Hash struct {
	message string
	hash [sha1.Size]byte

}

var counter int = 0
var max int = 1000000
var channel = make(chan Hash)

var producerWg sync.WaitGroup
var consumerWg sync.WaitGroup

func producer() {
	defer producerWg.Done()

	var source = rand.NewSource(time.Now().UnixNano())
	var generator = rand.New(source)

	for i := 0 ; i < max ; i++ {
		random := fmt.Sprintf("This is a test %v", generator.Int())
		hash := sha1.Sum([]byte(random))

		if (hash[0] == 0 && hash[1] == 0) {
			channel <- Hash{random, hash}
		}
	}
}

func consumer() {
	defer consumerWg.Done()

	for hash := range channel {
		fmt.Printf("Hash is %v\n ", hash)
    }
}

func main() {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	fmt.Println("Number of CPUs: ", nCPU)
	start := time.Now()

    rand.Seed(time.Now().Unix())

    for c := 0; c < runtime.NumCPU(); c++ {
        producerWg.Add(1)
        go producer()
    }

    for c := 0; c < runtime.NumCPU(); c++ {
        consumerWg.Add(1)
        go consumer()
    }

    producerWg.Wait()

    close(channel)

    consumerWg.Wait()

	fmt.Printf("Count of %v sha1 took %v\n", max, time.Since(start))
}


