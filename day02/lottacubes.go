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

var requiredCubes = map[string]int {
   "red": 12,
   "green": 13,
   "blue": 14,
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

func newGameSet(s string) map[string]int {
    m := make(map[string]int)

	tempNumColor := strings.SplitN(s, ",", -1)
	
    for i := 0; i < len(tempNumColor); i++ {
        TrimmedNumColor := strings.Trim(tempNumColor[i], " ")
        NumColor := strings.SplitN(TrimmedNumColor, " ", -1)

        m[NumColor[1]], _ = strconv.Atoi(NumColor[0])
    }

    return m
}
func CheckGame(mySet map[string]int, p *possibleGame) {

	if (mySet["red"] > requiredCubes["red"]) {
    	p.redIsPossible = false
	}
	if (mySet["green"] > requiredCubes["green"]) {
    	p.greenIsPossible = false
	}
	if (mySet["blue"] > requiredCubes["blue"]) {
    	p.blueIsPossible = false
	}
}

func PossibleGamesSum(s string, gn int, gt *int) {
    var isPossible possibleGame
    // Initialize everything to true. If everything was set to false,
    // we would have to check in both directions for every pass
    isPossible.redIsPossible 	= true
    isPossible.greenIsPossible 	= true
    isPossible.blueIsPossible 	= true
    // We will pass a pointer so we can go one map at a time and
    // still maintain the results
    var isPossiblePoint *possibleGame
    isPossiblePoint = &isPossible
    // We receive a string with the sets, not split
    // We proceed to split
    sets, numSets := splitSets(s)
    // We received a []string with sets and the number of sets
    // Now it's time to create a map
    var mySet map[string]int
    // For every set we have in the current game
    for i := 0; i < numSets; i++ {
        // We create a map
        mySet = newGameSet(sets[i])
        // We check if the game is possible
        CheckGame(mySet, isPossiblePoint)
    }

	if 	(isPossible.redIsPossible == true) &&
		(isPossible.greenIsPossible == true) &&
		(isPossible.blueIsPossible == true) {
    		*gt += gn
		} 
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
	// This variable will hold the pointer of the sum of possible games
	var gameSumPoint *int
	gameSum := 0
	gameSumPoint = &gameSum
	
	
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
    	PossibleGamesSum(gameAndCubes[1], gameNum, gameSumPoint)
	}
	fmt.Printf("The sum of possible games is: %d\n", gameSum) 
}
