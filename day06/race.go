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
    total		  	int
}

func (r *Race) WeRaceBoys(time, distance int, wg *sync.WaitGroup) {
    // As we just need to find the minimum necessary and then subtract it from
    // the maximum, probably a good idea starting form the middle
    tempTime := 0
    for i := int(time/2); i > 0; i-- {
        tempDist := i * (time - i)
        if tempDist <= distance {
            tempTime = i + 1
            break
        }
    }
    // If the minimum is tempTime, then the maximum is time - tempTime
    // time - 2*tempTime is the number of possible victories
    r.mu.Lock()
    r.total *= (time - (2 * tempTime - 1))
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

	file, err := os.Open("./inputs/day06_input")
	check(err)
	defer file.Close()

	r := Race{
		total: 1,
	}

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

	// E.g.: if I hold the button for 1ms and then release it, it will travel at
	// 1mm/ms fo the remaining amount of seconds.
	// We can skip the holding down 0 and tMAX ms.
	// Once we find the MIN amount of ms necessary to win, the limit is
	// MAX - MIN
	for i := 0; i < len(time); i++ {
    	wg.Add(1) 
    	go r.WeRaceBoys(time[i], distance[i], &wg) 
	}
	wg.Wait()
	PrintAndWait(r.total) 
}
