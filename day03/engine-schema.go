
package main

import (
    "fmt"
//    "strings"
    "bufio"
    "os"
    "regexp"
//    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func PrintAndWait(x ...any) {
    fmt.Print(x)
    fmt.Scanln() 
}

func main() {
	file, err := os.Open("./inputs/day03_input")
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
	numbers := make([][]int, numLines)
	_ = numbers

	/* 	Line [0] needs to check only line 1 for symbols.
		Similarly, the last line needs to check only
		(last line index - 1) for symbols.
		So we first check the line [0], then start a loop
		from 1 to (last line index - 1)
	*/
	firstLineNums 			:= renum.FindAllStringIndex(lines[0], -1)
	firstLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[0], -1)
	secondLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[1], -1)

	PrintAndWait(firstLineNums)
	PrintAndWait(firstLineSymbolsIndex) 
	PrintAndWait(secondLineSymbolsIndex)

	
	/* This code may need to be scrapped, not sure yet.
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
	*/
}
