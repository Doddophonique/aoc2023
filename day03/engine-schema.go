
package main

import (
    "fmt"
//    "strings"
    "bufio"
    "os"
    "regexp"
    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func PrintAndWait[T any](x T) {
    fmt.Print(x)
    fmt.Scanln() 
}

func main() {
	file, err := os.Open("./inputs/day03_test_input")
	check(err)
	defer file.Close()

	// This regex find the numbers inside strings
	renum := regexp.MustCompile("[0-9]+")
	resym := regexp.MustCompile("[^0-9.]+")

	scanner := bufio.NewScanner(file)
	var lines []string

	// Read the whole file into an array of strings
	for scanner.Scan() {
    	lines = append(lines, scanner.Text())
	}
	numLines := len(lines)

	// The 2D array of numbers will hold all the numbers to easily
	// match them with corresponding strings in the file using the index
	numbers := make([][]int, numLines)
	// For every line in the file, cerate an array of numbers
	for i := 0; i < numLines; i++ {
    	tempNums := renum.FindAllString(lines[i], -1)
    	for j := 0; j < len(tempNums); j++ {
    		num, _ := strconv.Atoi(tempNums[j])
    		numbers[i] = append(numbers[i], num)
    	}
	}

	// We store the index of a symbol on the line it appears
	symbolsIndex 	:= make([][]int, numLines)
	symbols 		:= make([][]string, numLines) 
	_ = symbolsIndex
	// For every line
	for i := 0; i < numLines; i++ {
		// We put all the symbols in a string
		tempSymbols := resym.FindAllString(lines[i], -1)
		// We associate symbols with an index
		tempSymbolsIndex := resym.FindAllStringIndex(lines[i], -1)
		// Add symbols and indexes to the arrays
		symbols = append(symbols, tempSymbols)
		//symbolsIndex = append(symbolsIndex, tempSymbolsIndex) 
	}
}
