package main

import (
    "fmt"
    "strings"
    "bufio"
    "math"
//    "math/rand"
    "os"
    "strconv"
//	"regexp"
	"sync"
)

var mapSeeds = map[string]int{
    "A": 13,
    "K": 12,
    "Q": 11,
    "J": 10,
    "T": 9,
    "9": 8,
    "8": 7,
    "7": 6,
    "6": 5,
    "5": 4,
    "4": 3,
    "3": 2,
    "2": 1,
}

type Game struct {
    mu 				sync.Mutex
    ranks		  	[]int
	// 0: High card, 1: One pair, 2: Two pair, 3: Three of a kind
	// 4: Full house, 5: Four of a kind, 6: Five of a kind
    typeOfHand		[7][]string
    indexOfHand		[7][]int
    baseThirteen    [7][]int
}

func (g *Game) ChangeBase(hType, index int, wg *sync.WaitGroup) {
	// Starting from the first char [0], we create the base13 num
	chars := len(g.typeOfHand[hType][index])
	baseTN := g.typeOfHand[hType][index]
	decNum := 0
	for i := 0; i < chars; i++ {
		// This should be refactored to be a bit more legible
		// It just computes N * 13^i and adds it over
    	decNum += int(float64(mapSeeds[string(baseTN[i])]) * math.Pow(13, float64(chars - i)))    }
	g.baseThirteen[hType][index] = decNum
    wg.Done()
}

func (g *Game) DetermineType(cards string, index int, wg *sync.WaitGroup) int {
    // We create a map and we check the length. Depending on the length, we
    // insert the string in a specific type
	m := make(map[string]int)
    for i := 0; i < len(cards); i++ {
        key := string(cards[i])
		m[key] = mapSeeds[key] 
    }
    // Now, depending on the number of elements in the map, we can assign
    // append cards to a specific rank
    mapSize := len(m)

    switch mapSize {
        // Five of a kind
        case 1:
           g.mu.Lock()
           g.typeOfHand[6] = append(g.typeOfHand[6], cards)
           g.indexOfHand[6] = append(g.indexOfHand[6], index) 
           g.mu.Unlock()
        // Four of a kind || Full House
        case 2:
           i := FullOrFour(cards)
           g.mu.Lock()
           g.typeOfHand[i] = append(g.typeOfHand[i], cards)
           g.indexOfHand[i] = append(g.indexOfHand[i], index) 
           g.mu.Unlock()
        // Three of a kind || Two pair
        case 3:
           i := ThreeOrTwo(cards)
           g.mu.Lock()
           g.typeOfHand[i] = append(g.typeOfHand[i], cards)
           g.indexOfHand[i] = append(g.indexOfHand[i], index) 
           g.mu.Unlock()
        // One pair
        case 4:
           g.mu.Lock()
           g.typeOfHand[1] = append(g.typeOfHand[1], cards)
           g.indexOfHand[1] = append(g.indexOfHand[1], index) 
           g.mu.Unlock()
        // High card
        case 5:
           g.mu.Lock()
           g.typeOfHand[0] = append(g.typeOfHand[0], cards)
           g.indexOfHand[0] = append(g.indexOfHand[0], index) 
           g.mu.Unlock()
    }
    wg.Done()

    return 1
}

func ThreeOrTwo(cards string) int {
    m := make(map[string]int)
    for i := 0; i < len(cards); i++ {
   		key := string(cards[i])
   		m[key] += 1
    }
    // m[i] returns 0 if the element is not in the map. I take advantage
    // of that
    for i := range mapSeeds {
        // If an element has 3 values, we have a three of a kind
        if m[i] == 3 {
            return 3
        /// If an element has 2 values, we have a two pair
        } else if m[i] == 2 {
            return 2
        }
    }
    return -1
}
func FullOrFour(cards string) int {
    m := make(map[string]int)
    for i := 0; i < len(cards); i++ {
   		key := string(cards[i])
   		m[key] += 1
    }
    // m[i] returns 0 if the element is not in the map. I take advantage
    // of that
    for i := range mapSeeds {
        // If an element has four values, we have a Four of a kind
        if m[i] == 4 {
            return 5
        /// If an element has 3 values, we have a Full House
        } else if m[i] == 3 {
            return 4
        }
    }
    return -1
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

// https://www.golangprograms.com/golang-program-for-implementation-of-quick-sort.html
// I need to learn how this shing works
func quicksort(a []int, g *Game, hType int) []int {
    if len(a) < 2 {
        return a
    }
     
    left, right := 0, len(a)-1
     
    pivot := 0
     
    a[pivot], a[right] = a[right], a[pivot]
    g.mu.Lock()
    g.indexOfHand[hType][pivot], g.indexOfHand[hType][right] = g.indexOfHand[hType][right], g.indexOfHand[hType][pivot]
    g.mu.Unlock()
     
    for i, _ := range a {
        if a[i] < a[right] {
            a[left], a[i] = a[i], a[left]
            g.mu.Lock()
            g.indexOfHand[hType][left], g.indexOfHand[hType][i] = g.indexOfHand[hType][i], g.indexOfHand[hType][left]
            g.mu.Unlock()
            left++
        }
    }
     
    a[left], a[right] = a[right], a[left]
    g.mu.Lock()
    g.indexOfHand[hType][left], g.indexOfHand[hType][right] = g.indexOfHand[hType][right], g.indexOfHand[hType][left]
    g.mu.Unlock()
    
     
    quicksort(a[:left], g, hType)
    quicksort(a[left+1:], g, hType)
     
    return a
}

func main() {

	file, err := os.Open("./inputs/day07_input")
	check(err)
	defer file.Close()

	// Struct for multiple races
	g := Game{}

	// Variable where we store every line in the file
	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
    	lines = append(lines, scanner.Text()) 
	}
	// Array of strings, set of cards
	var cards []string
	// Array of int, bet
	var bet []int
	// Now, split the lines
	for i := 0; i < len(lines); i++ {
    	tempString := strings.Split(lines[i], " ")
    	cards = append(cards, tempString[0])
    	tempNum, _ := strconv.Atoi(tempString[1])
    	bet = append(bet, tempNum) 
	}
	// Rank will be from 1 to len(lines)
	g.ranks = make([]int, len(lines))
	// What do we know for sure? 5 identical seeds are the highest ranks,
	// 5 completely different seeds are the lowest ranks.
	// We can iterate for every set of cards, and do different things
	// if the map we build has one element, five elements or the worst
	// case (two, three or four elements).
	//
	// Two identical: 	4oK or FH
	// Three identical: 3oK or 22
	// Four identical: 	12
	//
	// With this, I put every typeOfHand in a different array
	var wg sync.WaitGroup
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
    	g.DetermineType(cards[i], i, &wg)
	}

	wg.Wait()
	// Now that g.typeOfHand has every type of hand separated, the rank
	// is determined by the first card(s). I can convert every number
	// to a base 13 representation and sort.
	// In g.indexOfHand I have the index of the corresponding type.

	// For every type of hand
	for i := range g.typeOfHand {
    	// As many wait groups as the element we will iterate
    	wg.Add(len(g.typeOfHand[i]))
    	// We also need to initialize the array g.baseThirteen so we can
    	// keep the same index 
    	g.baseThirteen[i] = make([]int, len(g.typeOfHand[i]))
    	// For every element in a single type of hand
    	for j := range g.typeOfHand[i] {
			g.ChangeBase(i, j, &wg) 
    	}
	}
	wg.Wait()


	PrintAndWait(g.typeOfHand[0]) 
	PrintAndWait(g.typeOfHand[1]) 
	PrintAndWait(g.typeOfHand[2]) 
	PrintAndWait(g.typeOfHand[3]) 
	PrintAndWait(g.typeOfHand[4]) 
	PrintAndWait(g.typeOfHand[5]) 
	PrintAndWait(g.typeOfHand[6]) 
	// A sort of some kind. Important is to also move the index with the number as
	// well.
	//
	for i := range g.baseThirteen {
    	quicksort(g.baseThirteen[i], &g, i)    	
	}
	
	curRank := 1
	rank := 0
	// Iter every array
	for i := range g.typeOfHand {
    	for j := range g.typeOfHand[i] {
        	index := g.indexOfHand[i][j]
        	rank += curRank * bet[index]
        	curRank++
    	}
	}
	fmt.Print(rank) 
}
