package main

import (
    "strconv"
    "fmt"
)

/*
Multiples of 3 and 5
Problem 1
If we list all the natural numbers below 10 that are multiples of 3 or 5, we get 3, 5, 6 and 9. The sum of these multiples is 23.

Find the sum of all the multiples of 3 or 5 below 1000.
 */

func main() {
    max := 100000000
    picked := make([]int, 0)
    sum := 0

    for head := 0; head < max; head++ {
        str := strconv.FormatInt(int64(head), 10)
        last := str[len(str) - 1]
        if last == 0x30 || last == 0x35 {
            sum += head
            picked = append(picked, head)
            continue
        }
        sumDigits := 0
        for _, digit := range str {
            sumDigits += int(digit - 48)
        }
        if (sumDigits % 3 == 0) {
            sum += head
            picked = append(picked, head)
            continue
        }
    }

    fmt.Printf("Sum of natural numbers below %v that are multiplies of 3 or 5 is %v\n", max, sum)
//    fmt.Printf("Picked numbers are %v", picked)
}

