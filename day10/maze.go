package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

// Parallel code, global vars
type Nodes struct {
	mu       sync.Mutex
	variable int
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

func main() {
	file, err := os.Open("./inputs")
	check(err)
	defer file.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
}
