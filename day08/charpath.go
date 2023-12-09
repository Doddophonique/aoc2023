package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sync"
)

const LEFT = 'L'

type Nodes struct {
	mu       sync.Mutex
	commands []int32
	singleN  []int32
	leftN    []int32
	rightN   []int32
	index    int
	steps    uint64
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

func (n *Nodes) toByteSingle(s string) {
	// I've just received something like AAA
	var temp int32
	for i := len(s) - 1; i >= 0; i-- {
		a := int32(s[i])
		temp += a << ((len(s) - 1 - i) * 8)
	}
	n.singleN = append(n.singleN, temp)
}

func (n *Nodes) toByteDuet(s, r string) {
	// I've just received something like AAA BBB
	var tempL, tempR int32
	for i := len(s) - 1; i >= 0; i-- { 
		tempL += int32(s[i]) << ((len(s) - 1 - i) * 8)
		tempR += int32(r[i]) << ((len(s) - 1 - i) * 8)
		//fmt.Printf("Received %c (%b), now I have %b | ", s[i], s[i], tempL)
		//fmt.Printf("Received %c (%b), now I have %b\n", r[i], r[i], tempR)
		//fmt.Scanln() 
	}
	n.leftN = append(n.leftN, tempL)
	n.rightN = append(n.rightN, tempR)
}

func (n *Nodes) findNext(myN int32) {
	//var wg sync.WaitGroup
	for i := 0; i < len(n.singleN); i++ {
    	if myN^n.singleN[i] == 0 {
        	n.index = i
        	break
    	}
	}
}

func main() {
	file, err := os.Open("./inputs/day08_input")
	check(err)
	defer file.Close()
	// Struct with my node
	n := Nodes{}
	// Prepare the regex
	repath := regexp.MustCompile("([A-Z]{3})")
	// Build the END
	var END = 'Z'
	END += ('Z' << 8)
	END += ('Z' << 16)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	// First line, RL commands
	strCommands := scanner.Text()
	// Get every char inside the string just obtained
	for i := 0; i < len(strCommands); i++ {
		n.commands = append(n.commands, int32(strCommands[i]))
	}
	// One empty line
	scanner.Scan()
	// X = (Y, Z)
	// We regex this one
	for scanner.Scan() {
		tempNodes := repath.FindAllString(scanner.Text(), -1)
		n.toByteSingle(tempNodes[0])
		n.toByteDuet(tempNodes[1], tempNodes[2])
	}
	// We start from 0, we find the match
	// Let's start an infinite loop
	// Circular index
	i := 0
	// We start from the AAA element
	START := 'A'
	START += ('A' << 8)
	START += ('A' << 16)
	matching := START
	// Store where AAA is
	n.findNext(matching) 
	// By default, we will assume we are on the right
	// If we are not, we are in the left
	matching = n.rightN[n.index] 
	if n.commands[i]^LEFT == 0 {
		matching = n.leftN[n.index]
	}
	n.steps++
	// Infinite loop
	for {
		// Every step is in a single direction. For every step, we may need to
		// scan len(n.singleN) elements.
		// Circular loop
		n.findNext(matching)
		// Increment i after finding the match
		i++
		i = i % len(n.commands)
		// By default, we will assume we are on the right
		matching = n.rightN[n.index]
		//PrintAndWait()
		// If we are not, we are in the left
		if n.commands[i]^LEFT == 0 {
			matching = n.leftN[n.index]
		}
		n.steps++
		// If we find ZZZ, end
		if matching^END == 0 {
			break
		}
	}
	fmt.Printf("\nSteps: %d\n", n.steps)
}
