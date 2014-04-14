package main

import (
	"bufio"
	"os"
	"bytes"
	"strconv"
	"fmt"
	"time"
	"sort"
	"strings"
)

const (
	InputFilename       = "D-small-attempt0"
	ProblemURL          = "https://code.google.com/codejam/contest/2974486/dashboard#s=p3"
	CaseInputLinesCount = 3

	InputFilenameExtension  = ".in"
	OutputFileExtension = ".out"
)

var (
	TotalTimeCounter time.Duration = 0
	OutputBuffer bytes.Buffer
)

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

func solveCase(index int, input [][]float64) {
	//	log("Solving case %v, with data %v", index, input)

	naomi := input[1]
	ken := input[2]

	cheatNaomi := make([]float64, len(naomi))
	copy(cheatNaomi, naomi)
	sort.Float64s(cheatNaomi)
	cheatKen := make([]float64, len(ken))
	copy(cheatKen, ken)
	sort.Float64s(cheatKen)
	warNaomi := make([]float64, len(naomi))
	copy(warNaomi, naomi)
	sort.Float64s(warNaomi)
	warKen := make([]float64, len(ken))
	copy(warKen, ken)
	sort.Float64s(warKen)

	out("Case #%v: %v %v", index+1, cheat(cheatNaomi, cheatKen), war(warNaomi, warKen))
}

func cheat(naomi, ken []float64) int {
	log("Cheat game with Naomi's: %v and Ken's: %v", naomi, ken)
	score := 0
	rounds := len(naomi)

	for i := 0; i < rounds; i++ {
		//		maxKen := max(ken)
		minKen := min(ken, 0)

		minNaomi := min(naomi, minKen)
		shouldCheat := minNaomi > 0

		log("Min Ken is %v, min Naomi is %v, shouldCheat: %v", minKen, minNaomi, shouldCheat)

		if (shouldCheat) {
			score += 1
			for nidx, naomis := range naomi {
				if (naomis == minNaomi) {
					naomi[nidx] = -1
					break
				}
			}
		} else {
			for nidx, naomis := range naomi {
				if (naomis > -1) {
					naomi[nidx] = -1
					break
				}
			}
		}

		for kidx, kens := range ken {
			if (kens == minKen) {
				ken[kidx] = -1
				break;
			}
		}
	}
	log("Cheat game Naomi's score is: %v", score)
	return score

}

func war(naomi, ken []float64) int {
	log("War game with Naomi's: %v and Ken's: %v", naomi, ken)
	score := 0
	for _, naomis := range naomi {
		found := false
		for kidx, kens := range ken {
			if (kens > naomis) {
				ken[kidx] = -1
				found = true
				break;
			}
		}
		if (found == false) {
			for kidx, kens := range ken {
				if (kens > -1) {
					ken[kidx] = -1
					break;
				}
			}
			score += 1
		}
	}
	log("War game Naomi's score is: %v", score)
	return score
}

func max(floats []float64) float64 {
	max := -1.0
	for _, floaty := range floats {
		if (floaty > max) {
			max = floaty;
		}
	}
	return max
}

func min(floats []float64, lowest float64) float64 {
	min := 2.0
	for _, floaty := range floats {
		if (floaty >= lowest && floaty < min) {
			min = floaty;
		}
	}
	if min == 2.0 {
		min = -1.0
	}
	return min
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
		formatted = fmt.Sprint(format + "\n")
	} else {
		formatted = fmt.Sprintf(format + "\n", a...)
	}
	fmt.Print(formatted)
	OutputBuffer.WriteString(formatted)
}
