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

	// Go through each line to find matching runes and boxes
	correctBoxes := make([]string, 0)
	for _, line1 := range allLines {
		checkSum(line1, &twoSum, &threeSum)
		for _, line2 := range allLines {

			// We don't want to add the match of a line to itself
			if line1 == line2 {
				continue
			}

			// If the count of matching letters is 1 - len(line), add both to the correct box list
			matchCount := 0
			for i, _ := range line1 {

				// Count matching letters
				if line1[i] == line2[i] {
					matchCount++
				}

				// Add the boxes if we've found a ~match
				if matchCount == (len(line1) - 1) {
					correctBoxes = append(correctBoxes, line1)
					correctBoxes = append(correctBoxes, line2)
				}

			}

		}
	}

	// Get the letters in common between the correct boxes
	box1 := correctBoxes[0]
	box2 := correctBoxes[1]
	commonLetters := ""
	for i := range box1 {
		if box1[i] == box2[i] {
			commonLetters += string(box1[i])
		}
	}

	// Generate the check sum
	checkSum := twoSum * threeSum

	// Output the answer
	end := time.Now()
	fmt.Println()
	fmt.Println("Check Sum: ", checkSum)
	fmt.Println("Box Letters: ", commonLetters)
	fmt.Println("Time Taken: ", end.Sub(start))
	fmt.Println()

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
