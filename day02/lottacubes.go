package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type gameSet struct {
    nCubes []int
    color  []string    
}

type possibleGame struct {
    redIsPossible 	bool
    greenIsPossible bool
    blueIsPossible	bool
}

func splitSets (s string) ([]string, int) {
	// Every set is divided by ;
	sets := strings.SplitN(s, ";", -1)
	// Number of sets can vary, we need to have that info
	numSets := len(sets)

	return sets, numSets
}

func newGameSet(s []string, i int) []gameSet {
    // Make an array of structs as big as we need
    mySets := make([]gameSet, i)
    _ = mySets
	// Iterate the array of strings
	// e.g., s = 4 red, 5 blue, 7 green
	for index, set := range s {
    	// We extract every set
    	tempNumAndColor := strings.SplitN(set, ",", -1)
    	_ = index
    	mySets[index].nCubes = make([]int, len(tempNumAndColor))
    	mySets[index].color = make([]string, len(tempNumAndColor))
    	// The order is always number-space-color
    	// We now populate the struct
    	for j := 0; j < len(tempNumAndColor); j++ {
        	// Remove leading whitespace
        	TrimmedNumColor := strings.Trim(tempNumAndColor[j], " ")
        	// Split into number and color
        	NumColor := strings.SplitN(TrimmedNumColor, " ", -1)
        	// [0] is number, [1] is color
        	mySets[index].nCubes[j], _ = strconv.Atoi(NumColor[0])
        	mySets[index].color[j] = NumColor[1]
    	}
	}
	return mySets
}

func CheckGame(mySets []gameSet) bool {

    return false
}

func PossibleGamesSum(s string) bool {
    // We receive a string with the sets, not split
    // We proceed to split
    sets, numSets := splitSets(s)
    // We received a []string with sets and the number of sets
    // Now it's time to create a struct
    mySets := newGameSet(sets, numSets)
    isPossible := CheckGame(mySets)
	_ = isPossible
    return false
    
}

func PrintAndWait[T any](x T) {
    fmt.Print(x)
    fmt.Scanln() 
}

func main() {
	file, err := os.Open("input")
	check(err)
	defer file.Close()

	// This variable will hold the game number
	var gameNum int = 0
	_ = gameNum
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
    	line := scanner.Text()
    	// Split the string in "Game N" and "cubes color"
    	gameAndCubes := strings.Split(line, ":")
    	// At this point, gameAndCubes[0] has "Game N"
    	// We convert the number that remains after replacing "Game " with ""
    	gameNum, _ = strconv.Atoi(strings.Replace(gameAndCubes[0], "Game ", "", 1))
    	// Now, for every game, split the sets
    	PossibleGamesSum(gameAndCubes[1]) 
	}
}
