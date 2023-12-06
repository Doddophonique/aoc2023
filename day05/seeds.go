package main

import (
    "fmt"
//    "strings"
    "bufio"
    "math"
    "os"
    "strconv"
	"regexp"
)

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

func NextResource(previous []int, resource [][]int) []int {
	tempRes := make([]int, 0)
    for i := range previous {
        for j := range resource {
        	if previous[i] >= resource[j][1] &&
        	   previous[i] <= (resource[j][1] + resource[j][2] - 1) {
        		tempRes = append(tempRes, previous[i] + (resource[j][0] - resource[j][1]))
        	} 
        }
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

	// [0] is dest, [1] is source, [2] is range
	tempRes := make([]int, 0)
	for i := range seeds {
    	for j := range soils {
        	if seeds[i] >= soils[j][1] &&
        	   seeds[i] <= (soils[j][1] + soils[j][2] - 1) {
            	   tempRes = append(tempRes, (seeds[i] + (soils[j][0] - soils[j][1]))) 
        	}
    	}
    	if len(tempRes) == i {
        	tempRes = append(tempRes, seeds[i]) 
    	}
	}
	tempRes = NextResource(tempRes, fertilizers)
	tempRes = NextResource(tempRes, waters) 
	tempRes = NextResource(tempRes, lights) 
	tempRes = NextResource(tempRes, temperatures) 
	tempRes = NextResource(tempRes, humidities) 
	tempRes = NextResource(tempRes, locations)

	minimum := math.MaxInt
	for i := range tempRes {
    	if tempRes[i] < minimum {
        	minimum = tempRes[i] 
    	}
	}
	PrintAndWait(minimum) 
}
