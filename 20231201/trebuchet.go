package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "strconv"
)

var numbers = []string {
   "one", "two", "three", "four", "five",
   "six", "seven", "eight", "nine",
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ExtractCalibration(s string) string {

	firstDigit := 0
	lastDigit := 0
	index := 0

	// Start from the beginning of the string, stop when you find a number
    for i := 0; i < len(s); i++ {
       if num, err := strconv.Atoi(string(s[i])); err == nil {
           firstDigit = num
           index = i
           break
       }
    }

	// Start from the end of the string, stop when you find a number
    for j := len(s) - 1; j >= index; j-- {
       if num, err := strconv.Atoi(string(s[j])); err == nil {
           lastDigit = num
           break
       }
    }
   
    
    return fmt.Sprint(firstDigit) + fmt.Sprint(lastDigit)
}

func SearchFirstIndex(s string) int {
    // len(s) will always be bigger than the first index where you
    // can have a substring match
    index, number := len(s),0
    // from 1 to 9, save the smallest index and the number
    for i := 1; i < 10; i++ {
        tempIndex := strings.Index(s, fmt.Sprint(i))
        if tempIndex < index && tempIndex != -1 {
            index = tempIndex
            number = i
        }

    }
    // again for strings, from [0] "one" to [8] "nine"
    // save the smallest index, if smaller than the previous
    for i := 0; i < 9; i++ {
        tempIndex := strings.Index(s, numbers[i])
        if tempIndex < index && tempIndex != -1 {
            index = tempIndex
            number = i + 1
        }
    }

    return number
}

func SearchLastIndex(s string) int {
    // -1 will always be smaller than the last index where you
    // can have a substring match
    index, number := -1,0
    for i := 1; i < 10; i++ {
        tempIndex := strings.LastIndex(s, fmt.Sprint(i))
        if tempIndex > index && tempIndex != -1 {
            index = tempIndex
            number = i
        }

    }
    for i := 0; i < 9; i++ {
        tempIndex := strings.LastIndex(s, numbers[i])
        if tempIndex > index && tempIndex != -1 {
            index = tempIndex
            number = i + 1
        }
    }

    return number
}

func CombineIndexes(s string) int {
    firstNumber := SearchFirstIndex(s)
    secondNumber := SearchLastIndex(s)

    combinedString := fmt.Sprint(firstNumber) + fmt.Sprint(secondNumber)

    combinedNumber, err := strconv.Atoi(combinedString)
    _ = err

    return combinedNumber    
}

func main() {
	file, err := os.Open("input")
	check(err)
	defer file.Close()

	var firstTotal int = 0
	var secondTotal int = 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
    	line := scanner.Text()
    	calibration := ExtractCalibration(line)
    	num, err := strconv.Atoi(calibration)
    	_ = err
    	firstTotal += num
    	secondTotal += CombineIndexes(line) 
	}
	fmt.Printf("Calibration value, first method: %d\n", firstTotal)
	fmt.Printf("Calibration value, second method: %d\n", secondTotal) 
	
}
