package main

import (
    "crypto/sha1"
    "fmt"
    "time"
    "runtime"
    "math/rand"
)

type Hash struct {
    message string
    hash [sha1.Size]byte

}

var counter int = 0
var max int = 100000
var channel = make(chan Hash)
var source = rand.NewSource(time.Now().UnixNano())
var generator = rand.New(source)

func main() {
    nCPU := runtime.NumCPU()
    runtime.GOMAXPROCS(nCPU)
    fmt.Println("Number of CPUs: ", nCPU)
    start := time.Now()

    for i := 0 ; i < max ; i++ {
        go func(j int) {
            count(j)
        }(i)
    }

    for {
        select {
        case hash := <- channel:
            fmt.Printf("Hash is %v\n ", hash)
//		case <- done:
			break
        case <- time.After(50 * time.Millisecond):
            fmt.Printf(".")
        }
    }
    fmt.Printf("Count of %v sha1 took %v\n", max, time.Since(start))
}

func count(i int) {
    random := fmt.Sprintf("This is a test %v", generator.Int())
    hash := sha1.Sum([]byte(random))

    if (hash[0] == 0 && hash[1] == 0) {
        channel <- Hash{random, hash}
    }
}


