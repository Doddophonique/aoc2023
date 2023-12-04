package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
//    "strconv"
	"regexp"
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

func SplitSets(s []string) ([]string, []string) {

    // Split the two sets into winning and my numbers
	tempwinningNumbers := s[0]
	tempMyNumbers := s[1]
	// Regex to populate a string with numbers
	renum := regexp.MustCompile("[0-9]+")
	myNumbers := renum.FindAllString(tempMyNumbers, -1)
	winningNumbers := renum.FindAllString(tempwinningNumbers, -1)

	return myNumbers, winningNumbers  
}

func FindMatches(myNum []string, winNum []string) int {
    matches := 0 
	for i := range myNum {
    	for j := range winNum {
        	if (myNum[i] == winNum[j]) {
            	matches += 1
            	continue
        	}
    	}
	}
	return matches
}

func CalcTickets(s []string, index int, tpr []int ) {
    
	myNumbers, winningNumbers := SplitSets(s)
	matches := FindMatches(myNumbers, winningNumbers)
	if (matches > 0) {
    	for j := 0; j < tpr[index]; j++ {
        	for i := index; i < index + matches; i++ {
            	tpr[i+1] += 1
        	}
        	
    	}
	}
}

func CalcScore (s []string, t *int) {

	myNumbers, winningNumbers := SplitSets(s) 
	matches := FindMatches(myNumbers, winningNumbers)
	if(matches > 0) {
    	tempTotal := 1
    	for i := 0; i < matches - 1; i++ {
        	tempTotal *= 2
    	}
    	*t += tempTotal
	}
}

func main() {
	file, err := os.Open("./inputs/day04_input")
	check(err)
	defer file.Close()

	// The total is a simple sum 
	var total int = 0
	var totalPoint *int = &total
	// To keep track of the amount of tickets, an array as long
	// as the number of lines
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
    	lines = append(lines, scanner.Text())
	}
	ticketsPerRow := make([]int, len(lines))
	tprPoint := ticketsPerRow[0:len(lines)]

	for i := range ticketsPerRow {
    	ticketsPerRow[i] = 1
	}
	
	// Scan every line
	for i := 0; i < len(lines); i++ {
    	// e.g.: Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
    	cardAndNumbers := strings.Split(lines[i], ":")
    	// At this point, cardAndNumbers has "Card N"
    	// We don't need this information (yet?)
    	allNumbers := strings.Split(cardAndNumbers[1], "|")
    	CalcScore(allNumbers, totalPoint)
    	CalcTickets(allNumbers, i, tprPoint)
	}

	numTickets := 0
	for i := range ticketsPerRow {
    	numTickets += ticketsPerRow[i] 
	}
	fmt.Printf("The scratchcards are worth %d points.\n", total)
	fmt.Printf("In total, I have %d scratchcards.", numTickets)
}
