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
	var totalSum int = 0

	// Read the whole file into an array of strings
	for scanner.Scan() {
    	lines = append(lines, scanner.Text())
	}
	numLines := len(lines)
	numbers := make([][]int, numLines)

	// The 2D array of numbers will hold all the numbers to easily
	// match them with corresponding strings in the file using the index
	// For every line in the file, cerate an array of numbers
	for i := 0; i < numLines; i++ {
    	tempNums := renum.FindAllString(lines[i], -1)
    	for j := 0; j < len(tempNums); j++ {
    		num, _ := strconv.Atoi(tempNums[j])
    		numbers[i] = append(numbers[i], num)
    	}
	}
	/* 	Line [0] needs to check only line 1 for symbols.
		Similarly, the last line needs to check only
		(last line index - 1) for symbols.
		So we first check the line [0], then start a loop
		from 1 to (last line index - 1)
	*/
	firstLineNums 			:= renum.FindAllStringIndex(lines[0], -1)
	firstLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[0], -1)
	secondLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[1], -1)
	// For every *number index range*, check if in the same line there is a
	// symbol on (first - 1) or (last + 1), check in other lines if there is
	// a symbol in a specific interval of numbers. If you find a match, you
	// can break as you just need one symbol
	for i := range firstLineNums {
    	for j := range firstLineSymbolsIndex {
			if firstLineSymbolsIndex[j][0] >= firstLineNums[i][0] - 1 &&
			   (firstLineSymbolsIndex[j][0] <= firstLineNums[i][1]) {
    			   totalSum += numbers[0][i]
    			   break
			   }
    	}
    	for j := range secondLineSymbolsIndex {
			if (secondLineSymbolsIndex[j][0] >= firstLineNums[i][0] - 1)  &&
			   (secondLineSymbolsIndex[j][0] <= firstLineNums[i][1]) {
    			   totalSum += numbers[0][i]
    			   break
			   }
    	}
	}
	PrintAndWait(totalSum) 
	// Now we loop from 1 to (last index - i)
	for i := 1; i < len(lines) - 1; i++ {
    	// We need to check the current line against an interval of three lines
    	// breaking the loop for a single number as soon as we find a match
    	// (we don't want duplicate matches)
    	currentLineNums := renum.FindAllStringIndex(lines[i], -1)
    	previousLineIndex := resym.FindAllStringIndex(lines[i - 1], -1)
    	currentLineIndex := resym.FindAllStringIndex(lines[i], -1)
    	nextLineIndex := resym.FindAllStringIndex(lines[i + 1], -1)
    	PrintAndWait("i: ", i)
OuterLoop:
    	for k := range currentLineNums {
        	PrintAndWait("k: ", k, currentLineNums) 
    		for j := range previousLineIndex {
        		PrintAndWait("prev j: ", j, previousLineIndex) 
        		if previousLineIndex[j][0] >= currentLineNums[k][0] - 1 &&
        			previousLineIndex[j][0] <= currentLineNums[k][1] {
            			PrintAndWait(numbers[i][k]) 
            			totalSum += numbers[i][k]
            			continue OuterLoop
        			} 
    		}
    		for j := range currentLineIndex {
        		PrintAndWait("cur j: ", j, currentLineIndex) 
        		if currentLineIndex[j][0] >= currentLineNums[k][0] - 1 &&
        			currentLineIndex[j][0] <= currentLineNums[k][1] {
            			PrintAndWait(numbers[i][k]) 
            			totalSum += numbers[i][k]
            			continue OuterLoop
        			} 
    		}
    		for j := range nextLineIndex {
        		PrintAndWait("next j: ", j, nextLineIndex) 
        		if nextLineIndex[j][0] >= currentLineNums[k][0] - 1 &&
        			nextLineIndex[j][0] <= currentLineNums[k][1] {
            			PrintAndWait(numbers[i][k]) 
            			totalSum += numbers[i][k]
            			continue OuterLoop
        			} 
    		}
    	}
	}
	// Now we need to loop the last line and confront it with previous
	// and itself
	lastLineNums 			:= renum.FindAllStringIndex(lines[len(lines) - 1], -1)
	lastLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[len(lines) - 1], -1)
	notLastLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[len(lines) - 2], -1)
	// For every *number index range*, check if in the same line there is a
	// symbol on (last - 1) or (last + 1), check in other lines if there is
	// a symbol in a specific interval of numbers. If you find a match, you
	// can break as you just need one symbol
	for i := range lastLineNums {
    	PrintAndWait("i: ", i, lastLineNums) 
    	for j := range lastLineSymbolsIndex {
        	PrintAndWait("last j: ", j, notLastLineSymbolsIndex) 
			if lastLineSymbolsIndex[j][0] >= lastLineNums[i][0] - 1 &&
			   (lastLineSymbolsIndex[j][0] <= lastLineNums[i][1]) {
    			   PrintAndWait(numbers[len(lines) - 1][i])
    			   PrintAndWait(totalSum) 
    			   totalSum += numbers[len(lines) - 1][i]
    			   PrintAndWait(totalSum) 
    			   break
			   }
    	}
    	for j := range notLastLineSymbolsIndex {
        	PrintAndWait("notlast j: ", j, notLastLineSymbolsIndex) 
			if (notLastLineSymbolsIndex[j][0] >= lastLineNums[i][0] - 1) &&
			   (notLastLineSymbolsIndex[j][0] <= lastLineNums[i][1]) {
    			   PrintAndWait(numbers[len(lines) - 1][i])
    			   PrintAndWait(totalSum) 
    			   totalSum += numbers[len(lines) - 1][i]
    			   PrintAndWait(totalSum) 
    			   break
			   }
    	}
	}
	fmt.Printf("The total sum is: %d\n", totalSum)
}
