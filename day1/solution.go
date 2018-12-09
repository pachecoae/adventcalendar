package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Streamline error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// Track how long it takes to solve the problem
	start := time.Now()

	// Read the input file data
	data, err := os.Open("input.txt")
	check(err)

	// We need a variable to keep track of the total
	total := 0

	// We need a map to keep track of the total counts
	totalMap := make(map[int]int, 0)
	totalMap[total]++

	// Read the data line-by-line with a scanner
	// Each line is of the format '+#' or '-#'
	solutionFound := false
	fileScanner := bufio.NewScanner(data)
	allNumbers := make([]int, 0)
	for fileScanner.Scan() {

		// We need to convert each line into a string
		numString := fileScanner.Text()

		// Get the number (the sign is parsed by strconv.Atoi)
		num, err := strconv.Atoi(numString)
		check(err)

		// Add the number to the list
		allNumbers = append(allNumbers, num)

	}

	// We need to loop over after creating our total and search for a duplicate total value added to our map
	for !solutionFound {

		// Increment/Decrement the total and update the count in the map for each number
		for _, num := range allNumbers {
			total += num
			totalMap[total]++

			// Check for a duplicate total
			if totalMap[total] == 2 {
				fmt.Println("Total Count: ", total)
				solutionFound = true
				break
			}

		}

	}

	// Output the answer
	end := time.Now()
	fmt.Println()
	fmt.Println("Total Sum: ", total)
	fmt.Println("Time Taken: ", end.Sub(start))
	fmt.Println()

}
