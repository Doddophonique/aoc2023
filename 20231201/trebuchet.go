package main

import (
    "fmt"
    "bufio"
    "os"
    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ExtractCalibration(s string) string {

	firstDigit := 0
	lastDigit := 0
	index := 0
	
    for i := 0; i < len(s); i++ {
       if num, err := strconv.Atoi(string(s[i])); err == nil {
           firstDigit = num
           index = i
           break
       }
    }

    for j := len(s) - 1; j >= index; j-- {
       if num, err := strconv.Atoi(string(s[j])); err == nil {
           lastDigit = num
           break
       }
    }
   
    
    return fmt.Sprint(firstDigit) + fmt.Sprint(lastDigit)
}

func main() {
	file, err := os.Open("input")
	check(err)

	var total int = 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
    	line := scanner.Text()
    	calibration := ExtractCalibration(line)
    	num, err := strconv.Atoi(calibration)
    	_ = err
    	total += num
	}
	fmt.Println(total) 
}
