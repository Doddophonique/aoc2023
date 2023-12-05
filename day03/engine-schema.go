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

func MatchNumbers(lines []string, index int, ast []int, total *int) {
   	// Regex that finds the numbers in a row
	renum := regexp.MustCompile("[0-9]+")
    // Gather numbers in prev line
    prevLine := renum.FindAllStringIndex(lines[index-1], -1)
    // Gather numbers in this line
    thisLine := renum.FindAllStringIndex(lines[index], -1) 
    // Gather numbers in next line
    nextLine := renum.FindAllStringIndex(lines[index+1], -1)
    // Calculate the number of numbers in three lines
    totalNumbers := len(prevLine) + len(thisLine) + len(nextLine)

	// Now we create a big array with all the indexes
	allIndexes := prevLine
	for i := range thisLine {
    	allIndexes = append(allIndexes, thisLine[i])
	}
	for i := range nextLine {
    	allIndexes = append(allIndexes, nextLine[i]) 
	}

	// Now we create a big array with all the numbers
	// We start from the previous line
	allNumbers := renum.FindAllString(lines[index-1], -1)
    thisNums := renum.FindAllString(lines[index], -1)
    for i := range thisNums {
        allNumbers = append(allNumbers, thisNums[i]) 
    }
    nextNums := renum.FindAllString(lines[index+1], -1)
    for i := range nextNums {
        allNumbers = append(allNumbers, nextNums[i]) 
    }
	// When we start, we have zero matches
    matches := 0
    // We will stop when we encounter two numbers
    twoNums := [2]int{0, 0}
    _ = twoNums
    // Cycling through all numbers, but stopping at two matches
    for i := 0; i < totalNumbers && matches < 2; i++ {
        if (ast[0] >= allIndexes[i][0] - 1 && ast[0] <= allIndexes[i][1]) {
            matches += 1
            num, _ := strconv.Atoi(allNumbers[i])
            twoNums[matches - 1] = num
        }
    }
    if(matches == 2) {
        tempGears := twoNums[0] * twoNums[1]
        *total += tempGears
    }
}

func CheckGears(lines []string) {
 	total := 0
 	totalPoint := &total
   	// Regex that finds the numbers in a row
	//renum := regexp.MustCompile("[0-9]+")
	// Regex that finds the asterisks in a row
	resym := regexp.MustCompile("[*]")
	// For every line starting from the second
	for i := 0; i < len(lines) - 1; i++ {
    	// Take the index of the asterisks
    	asteriskIndex := resym.FindAllStringIndex(lines[i], - 1)
    	// For every index we get
    	for j := range asteriskIndex {
			MatchNumbers(lines, i, asteriskIndex[j], totalPoint) 
    	}    	
	}	
	// firstLineNums 			:= renum.FindAllStringIndex(lines[index], -1)
	// firstLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[index], -1)
	// secondLineSymbolsIndex 	:= resym.FindAllStringIndex(lines[index + 1], -1)
	// For every *number index range*, check if in the same line there is a
	// symbol on (first - 1) or (last + 1), check in other lines if there is
	// a symbol in a specific interval of numbers. If you find a match, you
	// can break as you just need one symbol
	fmt.Printf("Total of gears is: %d\n", total) 
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
	// For every line in the file, create an array of numbers
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
	// Now we loop from 1 to (last index - i)
	for i := 1; i < len(lines) - 1; i++ {
    	// We need to check the current line against an interval of three lines
    	// breaking the loop for a single number as soon as we find a match
    	// (we don't want duplicate matches)
    	currentLineNums := renum.FindAllStringIndex(lines[i], -1)
    	previousLineIndex := resym.FindAllStringIndex(lines[i - 1], -1)
    	currentLineIndex := resym.FindAllStringIndex(lines[i], -1)
    	nextLineIndex := resym.FindAllStringIndex(lines[i + 1], -1)
OuterLoop:
    	for k := range currentLineNums {
    		for j := range previousLineIndex {
        		if previousLineIndex[j][0] >= currentLineNums[k][0] - 1 &&
        			previousLineIndex[j][0] <= currentLineNums[k][1] {
            			totalSum += numbers[i][k]
            			continue OuterLoop
        			} 
    		}
    		for j := range currentLineIndex {
        		if currentLineIndex[j][0] >= currentLineNums[k][0] - 1 &&
        			currentLineIndex[j][0] <= currentLineNums[k][1] {
            			totalSum += numbers[i][k]
            			continue OuterLoop
        			} 
    		}
    		for j := range nextLineIndex {
        		if nextLineIndex[j][0] >= currentLineNums[k][0] - 1 &&
        			nextLineIndex[j][0] <= currentLineNums[k][1] {
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
    	for j := range lastLineSymbolsIndex {
			if lastLineSymbolsIndex[j][0] >= lastLineNums[i][0] - 1 &&
			   (lastLineSymbolsIndex[j][0] <= lastLineNums[i][1]) {
    			   totalSum += numbers[len(lines) - 1][i]
    			   break
			   }
    	}
    	for j := range notLastLineSymbolsIndex {
			if (notLastLineSymbolsIndex[j][0] >= lastLineNums[i][0] - 1) &&
			   (notLastLineSymbolsIndex[j][0] <= lastLineNums[i][1]) {
    			   totalSum += numbers[len(lines) - 1][i]
    			   break
			   }
    	}
	}
	fmt.Printf("The total sum is: %d\n", totalSum)
	CheckGears(lines) 
}
