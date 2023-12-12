package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Parallel code, global vars
type Maze struct {
	mu        sync.Mutex
	direction [][]rune
	steps     int
}

var mapDirections = map[rune]int8{
	'|': 0, // Vertical
	'-': 1, // Horizontal
	'L': 2, // North to East
	'J': 3, // North to West
	'7': 4, // South to West
	'F': 5, // South to East
	'.': 6, // Ground
	'S': 7, // Start
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func PrintAndWait(x ...any) {
	fmt.Print(x...)
	fmt.Scanln()
}

// use defer timer("funcname")() when the function you want to
// test starts
func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func (mz *Maze) ChooseFirstDirection(startPoint [2]int) [2]int {
	// For the first step, we don't care about directions, we only need a tile
	// were we can step on. E.g., if we have a J North of S we CAN'T go. A 7 would
	// be OK.
	x := startPoint[0]
	y := startPoint[1]
	// If x - 1 (going west), only 1, 2, 5 allowed
	temp := mapDirections[mz.direction[x-1][y]]
	if temp == 1 || temp == 2 || temp == 5 {
    	PrintAndWait("We go West boi")
    // If x + 1 (going east), only 1, 3, 4 allowed
	} else if temp := mapDirections[mz.direction[x+1][y]]; temp == 1 || temp == 3 || temp == 4 {
    	PrintAndWait("We go East boi")
  	// If y - 1 (going north), only 0, 4, 5 allowed
	} else if temp := mapDirections[mz.direction[x][y-1]]; temp == 0 || temp == 4 || temp == 5 {
    	PrintAndWait("We go North boi")
	} else {
    	log.Fatal("How is this even possible..?")
	}
}

func (mz *Maze) FindPath(startPoint [2]int) {
    nextCheck := startPoint
    nextCheck = mz.ChooseFirstDirection(nextCheck)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Must provide exactly one argument. It must be a file.")
		log.Fatal("NOW I WILL SCREAM AND DIE.")
	}
	file, err := os.Open(os.Args[1])
	check(err)
	defer file.Close()

	lines := make([]string, 0)

	mz := Maze{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	numColumns := len(lines)

	mz.direction = make([][]rune, numColumns)
	// Populate the array of directions
	for i := 0; i < numColumns; i++ {
		for j := 0; j < len(lines[i]); j++ {
			mz.direction[i] = append(mz.direction[i], rune(lines[i][j]))
		}
	}
	var startPoint [2]int
	// Search for the starting point
	for i := 0; i < numColumns; i++ {
		for j := 0; j < len(lines[i]); j++ {
			if mapDirections[mz.direction[i][j]] == 7 {
				startPoint[0] = i
				startPoint[1] = j
			}
		}
	}
	// Going west:  x - 1 | Going east:  x + 1
	// Going north: y - 1 | Going south: y + 1
	mz.FindPath(startPoint)
}
