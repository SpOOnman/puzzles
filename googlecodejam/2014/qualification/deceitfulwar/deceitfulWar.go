package main

import (
    "bufio"
    "os"
    "bytes"
    "strconv"
    "fmt"
    "time"
	"sort"
)

const Debug = false

//https://code.google.com/codejam/contest/351101/dashboard#s=p0
func main() {
    log("https://code.google.com/codejam/contest/351101/dashboard#s=p0")
    input, err := readFlatInput("/home/tkl/dev/puzzles/googlecodejam/2014/qualification/deceitfulwar/sample.in")
    if (err != nil) {
        log("Error %v", err)
        os.Exit(1)
    }

    casesCount := int(input[0][0])
    log("Solving %v cases", casesCount)
    for i := 0 ; i < casesCount ; i++ {
        start := time.Now()
        solveCase(i, input[1 + i * 3: 4 + i * 3])
        log("Took %v", time.Since(start))
    }
}

func solveCase(index int, input [][]float64) {
	log("Solving case %v, with data %v", index, input)

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

	out("Case #%v: %v %v", index + 1, cheat(cheatNaomi, cheatKen), war(warNaomi, warKen))
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
