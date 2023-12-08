package main

import(
	"fmt"
	"bufio"
	"os"
	"sync"
	"regexp"
	"context"
)

const LEFT = 'L'

type Nodes struct {
    mu 			sync.Mutex
    commands 	[]int32
    singleN		[]int32
    leftN		[]int32
    rightN		[]int32
    steps		uint64
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
    for i := 0; i < len(s); i++ {
        a := int32(s[i])
        temp +=  a << (i*8) 
    }
	n.singleN = append(n.singleN, temp)
}

func (n *Nodes) toByteDuet(s, r string) {
    // I've just received something like AAA BBB
    var tempL, tempR int32
    for i := 0; i < len(s); i++ {
        tempL += int32(s[i]) << (i*8)
        tempR += int32(s[i]) << (i*8)
    }
    n.leftN = append(n.leftN, tempL)
    n.rightN = append(n.rightN, tempR) 
}

func (n *Nodes) findNext(direction int32, index int, ctx context.Context) {
    
}

func main() {
	file, err := os.Open("./inputs/day08_test_input")
	check(err)
	defer file.Close()
	// Struct with my node
	n := Nodes{}
	// Prepare the regex
	repath := regexp.MustCompile("([A-Z]{3})")
	
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// We start from 0, we find the match
	for i, j := 0, 0; ; i++ {
    	// We do a circular loop
    	i =  i % len(n.commands)
    	// A function that has the context as an argument
    	n.findNext(int32(i), j, ctx)  
    	j++
	}
}
