package main

import (
    "fmt"
    "strings"
    "bufio"
//    "math"
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
           g.mu.Unlock()
        // Four of a kind || Full House
        case 2:
            //call func
           i := FullOrFour(cards)
           g.mu.Lock()
           g.typeOfHand[i] = append(g.typeOfHand[i], cards)
           g.mu.Unlock()
        // Three of a kind || Two pair
        case 3:
           i := ThreeOrTwo(cards)
           g.mu.Lock()
           g.typeOfHand[i] = append(g.typeOfHand[i], cards)
           g.mu.Unlock()
        // One pair
        case 4:
           g.mu.Lock()
           g.typeOfHand[1] = append(g.typeOfHand[1], cards)
           g.mu.Unlock()
        // High card
        case 5:
           g.mu.Lock()
           g.typeOfHand[0] = append(g.typeOfHand[0], cards)
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
            return 4
        /// If an element has 2 values, we have a two pair
        } else if m[i] == 2 {
            return 3
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
        /// If an element has 3 values, we have a four of a kind
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

func main() {

	file, err := os.Open("./inputs/day07_input")
	check(err)
	defer file.Close()

	// Struct for multiple races
	g := Game{}
	m := make(map[string]int)

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
	// Will try to call a different function for every line
	var wg sync.WaitGroup
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
    	go g.DetermineType(cards[i], i, &wg)
	}

	wg.Wait()
	PrintAndWait(g.typeOfHand) 
	_ = m
}
