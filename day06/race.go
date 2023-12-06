package main

import (
    "fmt"
//    "strings"
    "bufio"
//    "math"
    "os"
    "strconv"
	"regexp"
	"sync"
)

type Race struct {
    mu 				sync.Mutex
    temp_modify  	int
}

func (r *Race) WeRaceBoys(wg *sync.WaitGroup) {
    r.mu.Lock()
    r.mu.Unlock()
    wg.Done()
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

	file, err := os.Open("./inputs/day06_test_input")
	check(err)
	defer file.Close()

	r := Race{}
	_ = r

	var wg sync.WaitGroup
	_ = wg

   	// Regex that finds the numbers in a row
	renum := regexp.MustCompile("[0-9]+")
	_ = renum

	scanner := bufio.NewScanner(file)
	time, distance := make([]int, 0), make([]int, 0)
	tempStrings := make([]string, 0)
	// Generic, would work for more rows
	for scanner.Scan() {
    	tempStrings = append(tempStrings, scanner.Text())
	}
	// Now populate time and distance
	timeStr := renum.FindAllString(tempStrings[0], - 1)
	distStr := renum.FindAllString(tempStrings[1], - 1)

	// Both should be the same length
	for i := range timeStr {
    	num, _ := strconv.Atoi(timeStr[i])
    	time = append(time, num)
    	num, _ = strconv.Atoi(distStr[i])
    	distance = append(distance, num)
	}
	PrintAndWait(time, distance) 
}
