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

var fromWhere = map[rune]int8{
   'N': 0,
   'S': 1,
   'W': 2,
   'E': 3,
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

func (mz *Maze) ChooseFirstDirection(startPoint [2]int) ([2]int, rune) {
	// For the first step, we don't care about directions, we only need a tile
	// were we can step on. E.g., if we have a J North of S we CAN'T go. A 7 would
	// be OK.
	y := startPoint[0]
	x := startPoint[1]
	// If x - 1 (going west), only 1, 2, 5 allowed
	temp := mapDirections[mz.direction[y][x-1]]
	if temp == 1 || temp == 2 || temp == 5 {
    	// We came from East, we also return 'E'
    	newSpot := [2]int{y,x-1}
    	return newSpot, 'E'
    // If x + 1 (going east), only 1, 3, 4 allowed
	} else if temp := mapDirections[mz.direction[y][x+1]]; temp == 1 || temp == 3 || temp == 4 {
    	// We came from West
    	newSpot := [2]int{y,x+1}
    	return newSpot, 'W'
  	// If y - 1 (going north), only 0, 4, 5 allowed
	} else if temp := mapDirections[mz.direction[y-1][x]]; temp == 0 || temp == 4 || temp == 5 {
    	// We came from South
    	newSpot := [2]int{y-1,x}
    	return newSpot, 'S'
	} else {
    	log.Fatal("How is this even possible..?")
	}
	s := [2]int{0, 0}
	return s, 'X'
}

func (mz *Maze) WeGo(loc [2]int, direction rune) ([2]int, rune) {
    mapDir := fromWhere[direction]
    y, x := loc[0], loc[1]
    thisPipe := mapDirections[mz.direction[y][x]]
    goNorth := [2]int{y-1, x}
    goSouth := [2]int{y+1, x}
    goWest := [2]int{y, x-1}
    goEast := [2]int{y, x+1}
    switch mapDir {
        case 0:
            switch thisPipe {
                case 0:
                    return goSouth, 'N'
                case 2:
					return goEast, 'W'
                case 3:
                    return goWest, 'E'
            }
        case 1:
            switch thisPipe {
                case 0:
                    return goNorth, 'S'
                case 4:
                    return goWest, 'E'
                case 5:
                    return goEast, 'W'
            }
        case 2:
            switch thisPipe {
                case 1:
                    return goEast, 'W'
                case 3:
					return goNorth, 'S'
                case 4:
                    return goSouth, 'N'
            }
        case 3:
            switch thisPipe {
                case 1:
                    return goWest, 'E'
                case 2:
                    return goNorth, 'S' 
                case 5:
                    return goSouth, 'N'
            }
    }
	s := [2]int{0, 0}
	return s, 'X'
}

func (mz *Maze) FindPath(startPoint [2]int) {
    defer timer("FindPath")()
    y, x := startPoint[0], startPoint[1]
    nextCheck, direction := mz.ChooseFirstDirection(startPoint)
    mz.steps++
    curY, curX := nextCheck[0], nextCheck[1]
    for {
        nextCheck, direction = mz.WeGo(nextCheck, direction)
        curY, curX = nextCheck[0], nextCheck[1]
        mz.steps++
        if y == curY && x == curX { break } 
    }
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

	mz := Maze{ steps: 0, }

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
	fmt.Printf("Number of steps: %d", mz.steps / 2) 
}
