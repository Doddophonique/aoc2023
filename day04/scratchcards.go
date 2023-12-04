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

func CalcScore (s []string, t *int) {
    // Split the two sets into winning and my numbers
	tempwinningNumbers := s[0]
	tempMyNumbers := s[1]
	// Regex to populate a string with numbers
	renum := regexp.MustCompile("[0-9]+")
	myNumbers := renum.FindAllString(tempMyNumbers, -1)
	winningNumbers := renum.FindAllString(tempwinningNumbers, -1) 

	matches := 0
	for i := range myNumbers {
    	for j := range winningNumbers {
        	if (myNumbers[i] == winningNumbers[j]) {
            	matches += 1
            	continue
        	}
    	}
	}
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

	var total int = 0
	var totalPoint *int = &total

	scanner := bufio.NewScanner(file)

	// Scan every line
	for scanner.Scan() {
    	// e.g.: Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
    	line := scanner.Text()
    	cardAndNumbers := strings.Split(line, ":")
    	// At this point, cardAndNumbers has "Card N"
    	// We don't need this information (yet?)
    	allNumbers := strings.Split(cardAndNumbers[1], "|")
    	CalcScore(allNumbers, totalPoint)
	}
	fmt.Printf("The scratchcards are worth %d points.", total)
}
