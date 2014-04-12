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
    input, err := readFlatInput("/home/tkl/dev/puzzles/googlecodejam/2014/qualification/A-small-attempt0.in")
//    input, err := readFlatInput("/home/tkl/dev/puzzles/googlecodejam/2014/qualification/sample.in")
//    input, err := readFlatInput("/home/tkl/dev/googlecodejam/src/africa2010/A-large-practice.in")
    if (err != nil) {
        log("Error %v", err)
        os.Exit(1)
    }

    casesCount := input[0][0]
    log("Solving %v cases", casesCount)
    for i := 0 ; i < casesCount ; i++ {
        start := time.Now()
        solveCase(i, input[1 + i * 10 : 11 + i * 10])
        log("Took %v", time.Since(start))
    }
}

func solveCase(index int, input [][]int) {
    log("Solving case %v, with data %v", index, input)

	firstAnswer := input[0][0]
	firstRow := input[firstAnswer]
	secondAnswer := input[5][0]
	secondRow := input[5 + secondAnswer]

	log("First answer is %v row: %v, second answer is %v row: %v", firstAnswer, firstRow, secondAnswer, secondRow)
	candidates := make([]int, 0)

	for _, left := range firstRow {
		for _, right := range secondRow {
			if (left == right) {
				log("Found candidate: %v", left)
				candidates = append(candidates, left)
			}
		}
	}

	switch {
	case len(candidates) == 0:
		out("Case #%v: Volunteer cheated!", index + 1)
	case len(candidates) == 1:
		out("Case #%v: %v", index + 1, candidates[0])
	case len(candidates) > 1:
		out("Case #%v: Bad magician!", index + 1)
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
