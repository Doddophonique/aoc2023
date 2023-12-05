package main

import (
    "fmt"
//    "strings"
    "bufio"
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

func GetMaps(ss [][]string, nn [][]int) {

}

func main () {
	file, err := os.Open("./inputs/day05_test_input")
	check(err)
	defer file.Close()

   	// Regex that finds the numbers in a row
	renum := regexp.MustCompile("[0-9]+")
	
	var seeds []int
	//var soils, fertilizers, waters, lights, temperatures,
 	//	humidities, locations [][]int
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

}
