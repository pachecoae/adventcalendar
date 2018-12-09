package main

import (
	"bufio"
	"fmt"
	"os"
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

	// Read the data line-by-line with a scanner
	// Each line is of the format '+#' or '-#'
	fileScanner := bufio.NewScanner(data)
	allLines := make([]string, 0)
	for fileScanner.Scan() {

		// We need to convert each line into a string
		line := fileScanner.Text()

		// Add the line to the list
		allLines = append(allLines, line)

	}

	// Keep track of the values for the checksum
	twoSum := 0
	threeSum := 0

	// Go through each line to find matching runes
	for _, line := range allLines {
		checkSum(line, &twoSum, &threeSum)
	}

	// Generate the check sum
	checkSum := twoSum * threeSum
	fmt.Println(twoSum, threeSum)

	// Output the answer
	end := time.Now()
	fmt.Print(
		"\n",
		"Check Sum: ", checkSum,
		"\n",
		"Time Taken: ", end.Sub(start),
		"\n\n",
	)

}

func checkSum(line string, twoSum *int, threeSum *int) {

	// Find duplicated runes in the line
	runeMap := make(map[rune]int, 0)
	for _, r := range line {
		runeMap[r]++
	}

	// There can only be one counted for each
	two := false
	three := false

	// Check for the presence of 2 and 3 in the map
	for _, count := range runeMap {

		if !two && count == 2 {
			*twoSum++
			two = true
		}

		if !three && count == 3 {
			*threeSum++
			three = true
		}

	}

}
