package main

import (
	"bufio"
	"fmt"
	"math"
	"strings"
	"os"
	"strconv"
	"sync"
)

var mapSeedsFirst = map[string]int{
	"A": 12,
	"K": 11,
	"Q": 10,
	"J": 9,
	"T": 8,
	"9": 7,
	"8": 6,
	"7": 5,
	"6": 4,
	"5": 3,
	"4": 2,
	"3": 1,
	"2": 0,
}

var mapSeedsSecond = map[string]int{
	"A": 12,
	"K": 11,
	"Q": 10,
	"T": 9,
	"9": 8,
	"8": 7,
	"7": 6,
	"6": 5,
	"5": 4,
	"4": 3,
	"3": 2,
	"2": 1,
	"J": 0,
}

type Game struct {
	mu    sync.Mutex
	ranks []int
	// 0: High card, 1: One pair, 2: Two pair, 3: Three of a kind
	// 4: Full house, 5: Four of a kind, 6: Five of a kind
	typeOfHand   [7][]string
	indexOfHand  [7][]int
	baseThirteen [7][]int
}

func (g *Game) ChangeBase(hType, index int, mapSeeds map[string]int, wg *sync.WaitGroup) {
	// Starting from the first char [0], we create the base13 num
	chars := len(g.typeOfHand[hType][index])
	baseTN := g.typeOfHand[hType][index]
	decNum := 0
	for i := 0; i < chars; i++ {
		// This should be refactored to be a bit more legible
		// It just computes N * 13^i and adds it over
		decNum += int(float64(mapSeeds[string(baseTN[i])]) * math.Pow(13, float64(chars-i)))
	}
	g.baseThirteen[hType][index] = decNum
	wg.Done()
}

func (g *Game) AnalyzeMap(cards string, index, mapSize, offset int) {
	// The offset defaults to 0 for the regular game
	// it is +1 for a game with Jokers\
	mapSeed := mapSeedsFirst
	if offset != 0 {
		mapSeed = mapSeedsSecond
	}
	switch mapSize {
	// Five of a kind
	case 1:
		g.mu.Lock()
		g.typeOfHand[6] = append(g.typeOfHand[6], cards)
		g.indexOfHand[6] = append(g.indexOfHand[6], index)
		g.mu.Unlock()
	// Four of a kind || Full House
	case 2:
		i := FullOrFour(cards, offset, mapSeed)
		g.mu.Lock()
		g.typeOfHand[i] = append(g.typeOfHand[i], cards)
		g.indexOfHand[i] = append(g.indexOfHand[i], index)
		g.mu.Unlock()
	// Three of a kind || Two pair
	case 3:
		i := ThreeOrTwo(cards, offset, mapSeed)
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
}

func (g *Game) DetermineType(cards string, index int, wg *sync.WaitGroup) {
	// We create a map and we check the length. Depending on the length, we
	// insert the string in a specific type
	m := make(map[string]int)
	for i := 0; i < len(cards); i++ {
		key := string(cards[i])
		m[key] = mapSeedsFirst[key]
	}
	// Now, depending on the number of elements in the map, we can assign
	// append cards to a specific rank
	mapSize := len(m)
	g.AnalyzeMap(cards, index, mapSize, 0)

	wg.Done()
}

func (g *Game) SecondGame(cards string, index int, wg *sync.WaitGroup) {
	// If I don't have Js, run the standard DetermineType
	m := make(map[string]int)
	n := make(map[string]int)
	for i := 0; i < len(cards); i++ {
		key := string(cards[i])
		// We need to track the number of Js
		m[key] += 1
		// This is to track the type of hand, knowing the number of Js
		n[key] = mapSeedsSecond[key]
	}
	mapSize := len(n)
	switch m["J"] {
	case 0:
		// We have a hand without Js
		wg.Add(1)
		g.DetermineType(cards, index, wg)
	case 1:
		// If there is a J, J can be use is not adding information, therefore -1
		g.AnalyzeMap(cards, index, (mapSize - 1), 1)
	case 2:
		g.AnalyzeMap(cards, index, (mapSize - 1), 2)
	case 3:
		g.AnalyzeMap(cards, index, (mapSize - 1), 3)
	case 4:
		g.AnalyzeMap(cards, index, (mapSize - 1), 4)
	case 5:
		wg.Add(1)
		g.DetermineType(cards, index, wg)
	}
	wg.Done()
}

func ThreeOrTwo(cards string, offset int, mapSeed map[string]int) int {
	m := make(map[string]int)
	for i := 0; i < len(cards); i++ {
		key := string(cards[i])
		m[key] += 1
	}
	// If we are in the second game, remove J
	if mapSeed["J"] == 0 {
		m["J"] = 0
	}
	// m[i] returns 0 if the element is not in the map. I take advantage
	// of that
	tempNum := 0
	for i := range mapSeed {
		if m[i] > tempNum {
			tempNum = m[i]
		}
	}
	// If an element has 3 values, we have a three of a kind
	if tempNum+offset == 3 {
		return 3
		/// If an element has 2 values, we have a two pair
	} else if tempNum+offset == 2 {
		return 2
	}
	PrintAndWait("This has run for ", cards)
	return -1
}
func FullOrFour(cards string, offset int, mapSeed map[string]int) int {
	m := make(map[string]int)
	for i := 0; i < len(cards); i++ {
		key := string(cards[i])
		m[key] += 1
	}
	// If we are in the second game, remove J
	if mapSeed["J"] == 0 {
		m["J"] = 0
	}
	// m[i] returns 0 if the element is not in the map. I take advantage
	// of that
	// Better to save the maximum number
	tempNum := 0
	for i := range mapSeed {
		if m[i] > tempNum {
			tempNum = m[i]
		}
	}
	// If an element has four values, we have a Four of a kind
	if tempNum+offset == 4 {
		return 5
		/// If an element has 3 values, we have a Full House
	} else if tempNum+offset == 3 {
		return 4
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
func quicksort(a, h []int, hType int) []int {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := 0

	a[pivot], a[right] = a[right], a[pivot]
	h[pivot], h[right] = h[right], h[pivot]

	for i, _ := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			h[left], h[i] = h[i], h[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]
	h[left], h[right] = h[right], h[left]

	quicksort(a[:left], h[:left], hType)
	quicksort(a[left+1:], h[left+1:], hType)

	return a
}

func main() {

	file, err := os.Open("./inputs/day07_input")
	check(err)
	defer file.Close()

	// Struct for regular game
	g := Game{}
	// Struct for the Joker game
	jo := Game{}

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
	jo.ranks = make([]int, len(lines))
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
	wg.Add(2 * len(lines))
	for i := 0; i < len(lines); i++ {
		g.DetermineType(cards[i], i, &wg)
		jo.SecondGame(cards[i], i, &wg)
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
		wg.Add(len(jo.typeOfHand[i]))
		// We also need to initialize the array g.baseThirteen so we can
		// keep the same index
		g.baseThirteen[i] = make([]int, len(g.typeOfHand[i]))
		jo.baseThirteen[i] = make([]int, len(jo.typeOfHand[i]))
		// For every element in a single type of hand
		for j := range g.typeOfHand[i] {
			g.ChangeBase(i, j, mapSeedsFirst, &wg)
		}
		for j := range jo.typeOfHand[i] {
			jo.ChangeBase(i, j, mapSeedsSecond, &wg)
		}
	}
	wg.Wait()

	// A sort of some kind. Important is to also move the index with the number as
	// well.
	for i := range g.baseThirteen {
		quicksort(g.baseThirteen[i], g.indexOfHand[i], i)
		quicksort(jo.baseThirteen[i], jo.indexOfHand[i], i)
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
	fmt.Printf("Rank: %d\n", rank)

	curRank = 1
	rank = 0
	// Iter every array
	for i := range jo.typeOfHand {
		for j := range jo.typeOfHand[i] {
			index := jo.indexOfHand[i][j]
			rank += curRank * bet[index]
			curRank++
		}
	}
	fmt.Printf("Rank: %d\n", rank)

}
