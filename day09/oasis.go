package main

import(
    "fmt"
    "os"
    "bufio"
    "sync"
    "time"
    "regexp"
    "strconv"
)

// Parallel code, global vars
type Series struct {
	mu       	sync.Mutex
	numStore 	[][]int
	total		uint64
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

func PredictValueBack(numbers []int) []int {
	// Are we finished? By default, true
    temp := true
    for i := 0; i < len(numbers); i++ {
        if numbers[i] != 0 {
            temp = false
            break
        }
    }
    // Check to end recursion
    if temp == true {
        return numbers
    }

	newNums := make([]int, len(numbers) - 1)
	for i := 0; i < len(numbers) - 1; i++ {
    	newNums[i] = numbers[i + 1] - numbers[i]
	}
	myNums := PredictValueBack(newNums)
	addValue := newNums[0] - myNums[0]
	// We need to append at the start
	newNums = append([]int{addValue}, newNums...)
	return newNums
}

func (ser *Series) CallPredictBack(numbers []int, wg *sync.WaitGroup) {
    tempNum := PredictValueBack(numbers)
 	ser.mu.Lock()
    ser.total += uint64(numbers[0] - tempNum[0])
    ser.mu.Unlock()
    wg.Done()
}

func PredictValue(numbers []int) []int {
	// Are we finished? By default, true
    temp := true
    for i := 0; i < len(numbers); i++ {
        if numbers[i] != 0 {
            temp = false
            break
        }
    }
    // Check to end recursion
    if temp == true {
        return numbers
    }

	newNums := make([]int, len(numbers) - 1)	
	for i := 0; i < len(numbers) - 1; i++ {
    	newNums[i] = numbers[i + 1] - numbers[i]
	}
	myNums := PredictValue(newNums)
	addValue := myNums[len(myNums) - 1]
	newNums = append(newNums, newNums[len(newNums) - 1] + addValue)
	return newNums
}

func (ser *Series) CallPredict(numbers []int, wg *sync.WaitGroup) {
    tempNum := PredictValue(numbers)
    lt 		:= len(tempNum) - 1
 	ln		:= len(numbers) - 1
 	ser.mu.Lock()
    ser.total += uint64(numbers[ln] + tempNum[lt])
    ser.mu.Unlock()
    wg.Done()
}

func main() {
    defer timer("main")()
	file, err := os.Open("./inputs/day09_input")
	check(err)
	defer file.Close()

	var wg sync.WaitGroup

	renum := regexp.MustCompile("(\\-[0-9]+|[0-9]+)")

	ser := Series{ total: 0, }
	
	lines := make([]string, 0)
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	    lines = append(lines, scanner.Text())
	}
	ser.numStore = make([][]int, len(lines))
	for i := 0; i < len(lines); i++ {
    	temp := renum.FindAllString(lines[i], -1) 
    	for j := 0; j < len(temp); j++ {
        	num, err := strconv.Atoi(temp[j])
			check(err)
			ser.numStore[i] = append(ser.numStore[i], num)
    	}
	}
	// Now I have a 2D array with all the numbers, I can start RECURSING
	wg.Add(len(ser.numStore)) 
	for i := 0; i < len(ser.numStore); i++ {
    	go ser.CallPredict(ser.numStore[i], &wg)
	}
	wg.Wait()
	fmt.Printf("%d\n", ser.total)

	ser.total = 0
	wg.Add(len(ser.numStore))
	for i := 0; i < len(ser.numStore); i++ {
    	go ser.CallPredictBack(ser.numStore[i], &wg)
	}
	wg.Wait()
	fmt.Printf("%d\n", ser.total)
	
}   	
