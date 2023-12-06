package main

import (
    "fmt"
//    "strings"
    "bufio"
    "math"
    "os"
    "strconv"
	"regexp"
	"sync"
)

type Minimum struct {
    mu 		sync.Mutex
    minimum int
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

func GetMaps(scanner *bufio.Scanner, re *regexp.Regexp) [][]int {
 	// Scan until there is an empty line
 	var tempNums int = 0
 	tempArray := make([][]int, 0)
 	for i := 0; scanner.Scan() && scanner.Text() != ""; i++ {
     	tempString := re.FindAllString(scanner.Text(), -1)
     	temp := make([]int, 0)
     	for j := range tempString {
			tempNums, _ = strconv.Atoi(tempString[j])
			temp = append(temp, tempNums)
     	}
     	tempArray = append(tempArray, temp) 
 	}
 	// Prepare for next line
 	scanner.Scan()
 	return tempArray
}

func (min *Minimum) ParallelMinimum(start, finish int, atrocity [][][]int, wg *sync.WaitGroup) {
    // We check the array one by one, need a temp array because
    // SeedToLocation wants it
    tempNum := make([]int, 1)
    tempMin := []int{0}
    _ = tempMin
    for i := 0; i < finish; i++ {
        tempNum[0] = start + i
        tempMin := SeedToLocation(tempNum, atrocity)
        // Need to modify a shared variable, lock
        min.mu.Lock()
        if tempMin[0] < min.minimum {
            min.minimum = tempMin[0] 
        }
        min.mu.Unlock()
    }
    // We finished with the Goroutine
    wg.Done()
}

func SeedToLocation(seeds []int, atrocity [][][]int) []int {
	tempRes := seeds
	for i := range atrocity {
    	tempRes = NextResource(tempRes, atrocity[i])
	}
	return tempRes
}

func NextResource(previous []int, resource [][]int) []int {
	tempRes := make([]int, 0)
	// [0] is dest, [1] is source, [2] is range
    for i := range previous {
        for j := range resource {
        	if previous[i] >= resource[j][1] &&
        	   previous[i] <= (resource[j][1] + resource[j][2] - 1) {
        		tempRes = append(tempRes, previous[i] + (resource[j][0] - resource[j][1]))
        	} 
        }
        // If we didn't add an element to the array
    	if len(tempRes) == i {
        	tempRes = append(tempRes, previous[i])
    	}
    }
    return tempRes 
}
func main () {
	file, err := os.Open("./inputs/day05_input")
	check(err)
	defer file.Close()

	min := Minimum{
    	minimum: math.MaxInt,
	}

	var wg sync.WaitGroup

   	// Regex that finds the numbers in a row
	renum := regexp.MustCompile("[0-9]+")
	
	var seeds []int
	var soils,  fertilizers, waters, lights, temperatures,
 		humidities, locations [][]int
	scanner := bufio.NewScanner(file)
 	// We know that the seeds only have one row
 	scanner.Scan()
 	// Put all the numbers in an array of strings
 	seedsNums := renum.FindAllString(scanner.Text(), -1)
 	// Extract every number from the string
 	for i := 0; i < len(seedsNums); i++ {
		num, _ := strconv.Atoi(seedsNums[i])
		seeds = append(seeds, num) 
 	}
 	// We know we have an empty string and just a title, skip them
 	scanner.Scan()
 	scanner.Scan()
 	// Should be possible to just pass the scanner
 	soils = GetMaps(scanner, renum)
 	fertilizers = GetMaps(scanner, renum)
 	waters = GetMaps(scanner, renum)
 	lights = GetMaps(scanner, renum)
 	temperatures = GetMaps(scanner, renum)
 	humidities = GetMaps(scanner, renum)
 	locations = GetMaps(scanner, renum)

	tempRes := make([]int, 0)
	// Actually insane behaviour
	monster := [][][]int{
    	soils, fertilizers, waters,
    	lights, temperatures, humidities,
    	locations,
	}
	// Send the seeds, receive 
	tempRes = SeedToLocation(seeds, monster) 

	minimum := math.MaxInt
	for i := range tempRes {
    	if tempRes[i] < minimum {
        	minimum = tempRes[i] 
    	}
	}
	fmt.Printf("Minimum of first part: %d\n", minimum)

	// Actual madness
	for i := 0; i < len(seeds); i += 2 {
    	wg.Add(1)
    	go min.ParallelMinimum(seeds[i], seeds[i+1], monster, &wg) 
	}	

	wg.Wait()
	//tempRes = SeedToLocation(allSeeds, monster) 
	fmt.Printf("Minimum of second part: %d\n", min.minimum)
}
