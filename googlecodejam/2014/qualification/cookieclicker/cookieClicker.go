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
//    input, err := readFlatInput("/home/tkl/dev/puzzles/googlecodejam/2014/qualification/cookieclicker/sample.in")
    input, err := readFlatInput("/home/tkl/dev/puzzles/googlecodejam/2014/qualification/cookieclicker/B-small-attempt0.in")
    if (err != nil) {
        log("Error %v", err)
        os.Exit(1)
    }

    casesCount := int(input[0][0])
    log("Solving %v cases", casesCount)
    for i := 0 ; i < casesCount ; i++ {
        start := time.Now()
        solveCase(i, input[1 + i : 2 + i])
        log("Took %v", time.Since(start))
    }
}

func solveCase(index int, input [][]float64) {
    log("Solving case %v, with data %v", index, input)

	farmCost := input[0][0]
	farmGain := input[0][1]
	goal := input[0][2]

	var production float64 = 2

	result := simulate(production, farmCost, farmGain, goal, 0.0)
	out("Case #%v: %.7f", index + 1, result)
}

func simulate(production float64, farmCost float64, farmGain float64, goal float64, timeTaken float64) float64 {
	log("Current production is %v with total time taken %v", production, timeTaken)

	timeToWin := seconds(production, goal)
	timeToBuyNextFarm := seconds(production, farmCost)
	timeToWinWithNextFarm := timeToBuyNextFarm + seconds(production + farmGain, goal)

	shouldBuyAnotherFarm := timeToWin > timeToWinWithNextFarm

	log("Time to win: %v", timeToWin)
	log("Time to win with next farm: %v", timeToWinWithNextFarm)
	log("Should buy another farm: %v", shouldBuyAnotherFarm)

	if (shouldBuyAnotherFarm) {
		timeTaken += timeToBuyNextFarm
		return simulate(production + farmGain, farmCost, farmGain, goal, timeTaken)
	} else {
		return timeTaken + timeToWin
	}
}

func seconds(production float64, expected float64) float64 {
	return expected/production
}

//http://stackoverflow.com/questions/9862443/golang-is-there-a-better-way-read-a-file-of-integers-into-an-array
func readFlatInput(path string) ([][]float64, error) {
    file, err := os.Open(path)
    if (err != nil) {
        log("Error opening file: %v", err)
        os.Exit(1)
    }
    scanner := bufio.NewScanner(file)
    result  := make([][]float64, 0)
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
