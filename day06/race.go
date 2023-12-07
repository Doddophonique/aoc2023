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

func MultRaceDist(time, dist []string) ([]int, []int){
	tempT, tempD := make([]int, 0), make([]int, 0)
   	for i := range time {
    	num, _ := strconv.Atoi(time[i])
    	tempT = append(tempT, num)
    	num, _ = strconv.Atoi(dist[i])
    	tempD = append(tempD, num)
	}
    return tempT, tempD
}

func SingleRaceDist(time, dist []string) (int, int) {
    strT, strD := "", ""
	numT, numD := 0, 0
	// Create two big strings
   	for i := range time {
		strT += time[i]
		strD += dist[i]
   	}
   	
   	numT, _ = strconv.Atoi(strT) 
   	numD, _ = strconv.Atoi(strD) 
   		
    return numT, numD
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

	// Struct for multiple races
	rMult := Race{
		total: 1,
	}
	// Struct for a single race
	rSing := Race{
    	total: 1,
	}
 	_ = &rSing
	var wg sync.WaitGroup

   	// Regex that finds the numbers in a row
	renum := regexp.MustCompile("[0-9]+")

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

	time, distance = MultRaceDist(timeStr, distStr)

	// E.g.: if I hold the button for 1ms and then release it, it will travel at
	// 1mm/ms fo the remaining amount of seconds.
	// We can skip the holding down 0 and tMAX ms.
	// Once we find the MIN amount of ms necessary to win, the limit is
	// MAX - MIN
   	wg.Add(len(time)) 
	for i := 0; i < len(time); i++ {
    	go rMult.WeRaceBoys(time[i], distance[i], &wg) 
	}

	// Silly implementation of the single race
	singT, singD := SingleRaceDist(timeStr, distStr) 
	wg.Add(1)
	go rSing.WeRaceBoys(singT, singD, &wg) 
	
	wg.Wait()
	fmt.Printf("Multiple races result: %d.\n", rMult.total)
	fmt.Printf("Single race result: %d.\n", rSing.total) 
}
