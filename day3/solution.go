package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Struct to hold the parsed information
type rectangle struct {
	id     int
	left   int
	top    int
	width  int
	height int
}

// Take in a regular expression with text and return the matched text as a rectangle
func getRegexRectMatch(line string) (rectangle, error) {

	r, err := regexp.Compile("#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)$")
	check(err)

	allSubmatches := r.FindAllStringSubmatch(line, -1)
	i, l, t, w, h := "0", "0", "0", "0", "0"
	i, l, t, w, h =
		allSubmatches[0][1], allSubmatches[0][2], allSubmatches[0][3], allSubmatches[0][4], allSubmatches[0][5]

	id, err := strconv.Atoi(i)
	left, err := strconv.Atoi(l)
	top, err := strconv.Atoi(t)
	width, err := strconv.Atoi(w)
	height, err := strconv.Atoi(h)
	if err != nil {
		return rectangle{}, err
	}
	return rectangle{id, left, top, width, height}, nil

}

// Streamline error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Parse the input for the rectangle's distance from the left, distance from the top, width, and height
func parseLine(line string) (rectangle, error) {

	rect, err := getRegexRectMatch(line)
	if err != nil {
		return rectangle{}, err
	}
	return rect, nil

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
	allRectangles := make([]rectangle, 0)
	for fileScanner.Scan() {

		// We need to convert each line into a rectangle
		line := fileScanner.Text()
		rect, err := parseLine(line)
		if err != nil {
			continue
		}

		// Add the rectangle to the list
		allRectangles = append(allRectangles, rect)

	}

	// Keep track of overlapped area using a map of x -> y -> count
	rowMap := make(map[int]map[int]int, 0)
	for _, rect := range allRectangles {

		// Input data into the map for each (x, y) coordinate
		for x := rect.left; x < rect.left+rect.width; x++ {

			// Check to see if we need to initialize a column map
			if _, ok := rowMap[x]; !ok {
				rowMap[x] = make(map[int]int, 0)
			}

			// Counts will be recorded so that we can keep track of overlapping areas (count > 1)
			for y := rect.top; y < rect.top+rect.height; y++ {
				rowMap[x][y]++
			}

		}

	}

	// Run through the map a second time to find areas that have no overlaps
	id := 0
	for _, rect := range allRectangles {

		// Keep track of a rectangle with no collisions
		noCollisions := true

		// Input data into the map for each (x, y) coordinate
		for x := rect.left; x < rect.left+rect.width; x++ {

			// Check to see if there's a rectangle that has only 1 in each position in the map
			for y := rect.top; y < rect.top+rect.height; y++ {
				if rowMap[x][y] > 1 {
					noCollisions = false
				}
			}

		}

		if noCollisions {
			id = rect.id
		}

	}

	// Count overlapping areas
	overlappingAreas := 0
	for _, colMap := range rowMap {
		for _, count := range colMap {
			if count > 1 {
				overlappingAreas++
			}
		}
	}

	// Output the answer
	end := time.Now()
	fmt.Println()
	fmt.Println("Time Taken: ", end.Sub(start))
	fmt.Println("Overlapping Areas: ", overlappingAreas)
	fmt.Println("No Overlap Id: ", id)
	fmt.Println()

}
