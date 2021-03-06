package main

import (
	"bufio"
	"os"
	"bytes"
	"strconv"
	"fmt"
	"time"
	"strings"
)

const (
	InputFilename       = "A-small-attempt0"
	ProblemURL          = "https://code.google.com/codejam/contest/2974486/dashboard#s=p0"
	CaseInputLinesCount = 10

	InputFilenameExtension = ".in"
	OutputFileExtension    = ".out"
)

var (
	TotalTimeCounter time.Duration = 0
	OutputBuffer bytes.Buffer
)


func solveCase(index int, input [][]float64) {
    log("Solving case %v, with data %v", index, input)

	firstAnswer := int(input[0][0])
	firstRow := input[firstAnswer]
	secondAnswer := int(input[5][0])
	secondRow := input[5 + secondAnswer]

	log("First answer is %v row: %v, second answer is %v row: %v", firstAnswer, firstRow, secondAnswer, secondRow)
	candidates := make([]float64, 0)

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



func main() {
	log("Problem URL: %v", ProblemURL)
	input, err := readFlatInput()
	if (err != nil) {
		log("Error %v", err)
		os.Exit(1)
	}

	casesCount := int(input[0][0])
	log("Solving %v cases", casesCount)
	for i := 0 ; i < casesCount ; i++ {
		caseData := input[1+i*CaseInputLinesCount: (1+CaseInputLinesCount)+i*CaseInputLinesCount]
		log("Case #%v data: %v", i+1, caseData)
		start := time.Now()
		solveCase(i, caseData)
		duration := time.Since(start)
		TotalTimeCounter += duration
		log("Case #%v took %v, total time is %v\n\n", i+1, duration, TotalTimeCounter)
	}
	writeOutFile()
}

//http://stackoverflow.com/questions/9862443/golang-is-there-a-better-way-read-a-file-of-integers-into-an-array
func readFlatInput() ([][]float64, error) {
	wd, err := os.Getwd()
	file, err := os.Open(wd + "/" + InputFilename + InputFilenameExtension)
	if (err != nil) {
		log("Error opening file: %v", err)
		panic(1)
	}
	scanner := bufio.NewScanner(file)
	result := make([][]float64, 0)
	for scanner.Scan() {
		words := bytes.Split(scanner.Bytes(), []byte{0x20})
		float64s := make([]float64, 0, len(words))

		for _, word := range words {
			float64eger, err := strconv.ParseFloat(string(word), 64)
			if err != nil { return nil, err }
			float64s = append(float64s, float64eger)
		}
		result = append(result, float64s)
	}
	return result, nil
}

func writeOutFile() {
	wd, err := os.Getwd()
	file, err := os.Create(wd + "/" + InputFilename + OutputFileExtension)
	if (err != nil) {
		log("Error creating file: %v", err)
		panic(1)
	}
	_, err = file.WriteString(strings.TrimSpace(OutputBuffer.String()))
	if err = file.Close(); err != nil {
		panic(err)
	}
}

func log(format string, a ...interface{}) {
	if (a == nil) {
		fmt.Print(format + "\n")
	} else {
		fmt.Printf(format+"\n", a...)
	}
}

func out(format string, a ...interface{}) {
	var formatted string
	if (a == nil) {
		formatted = fmt.Sprint(format+"\n")
	} else {
		formatted = fmt.Sprintf(format+"\n", a...)
	}
	fmt.Print(formatted)
	OutputBuffer.WriteString(formatted)
}
