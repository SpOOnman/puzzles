//package africa2010
package main

import (
    "bufio"
    "os"
    "bytes"
    "strconv"
    "fmt"
    "time"
)

const Debug = false

//https://code.google.com/codejam/contest/351101/dashboard#s=p0
func main() {
    log("https://code.google.com/codejam/contest/351101/dashboard#s=p0")
    input, err := readFlatInput("/home/tkl/dev/googlecodejam/src/africa2010/A-small-practice.in")
//    input, err := readFlatInput("/home/tkl/dev/googlecodejam/src/africa2010/A-large-practice.in")
    if (err != nil) {
        log("Error %v", err)
        os.Exit(1)
    }

    casesCount := input[0][0]
    log("Solving %v cases", casesCount)
    for i := 0 ; i < casesCount ; i++ {
        start := time.Now()
        solveCase(i, input[1 + i * 3 : 4 + i * 3])
        log("Took %v", time.Since(start))
    }
}

func solveCase(index int, input [][]int) {
    log("Solving case %v, with data %v", index, input)

    credit := input[0][0]
    itemsCount := input[1][0]
    items := input[2]
    log("Credit is %v, there are %v items: %v", credit, itemsCount, items)

    for leftIndex, left := range items {
        remaining := items[leftIndex + 1:]
        for rightIndex, right := range remaining {
            if (left + right == credit) {
                out("Case #%v: %v %v", index + 1, leftIndex + 1, rightIndex + leftIndex + 1 + 1)
                l := items[leftIndex]
                r := items[rightIndex + leftIndex + 1]
                log("Verification: indices %v, %v cost %v, %v, together %v vs credit %v is %v", leftIndex, rightIndex + leftIndex + 1, l, r, l + r, credit, l + r == credit)
            }
        }
    }
}

//http://stackoverflow.com/questions/9862443/golang-is-there-a-better-way-read-a-file-of-integers-into-an-array
func readFlatInput(path string) ([][]int, error) {
    file, err := os.Open(path)
    if (err != nil) {
        log("Error opening file: %v", err)
        os.Exit(1)
    }
    scanner := bufio.NewScanner(file)
    result  := make([][]int, 0)
    for scanner.Scan() {
        words := bytes.Split(scanner.Bytes(), []byte{0x20})
        ints := make([]int, 0, len(words))

        for _, word := range words {
            integer, err := strconv.Atoi(string(word))
            if err != nil { return nil, err }
            ints = append(ints, integer)
        }
        result = append(result, ints)
    }
    return result, nil
}

func log(format string, a ...interface{}) (n int, err error) {
    if !Debug {
        return
    }
    if (a == nil) {
        return fmt.Print(format + "\n")
    } else {
        return fmt.Printf(format + "\n", a...)
    }
}

func out(format string, a ...interface{}) (n int, err error) {
    if (a == nil) {
        return fmt.Print(format + "\n")
    } else {
        return fmt.Printf(format + "\n", a...)
    }
}

func assert(condition bool) {

}
