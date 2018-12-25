package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// LogEntry to hold the parsed information
type LogEntry struct {
	date time.Time
	log  string
}

// By used to sort log entries
type By func(e1, e2 *LogEntry) bool

// Sort log entries
func (by By) Sort(logEntries []LogEntry) {
	les := &logEntrySorter{
		entries: logEntries,
		by:      by,
	}
	sort.Sort(les)
}

func (s *logEntrySorter) Len() int {
	return len(s.entries)
}

func (s *logEntrySorter) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

func (s *logEntrySorter) Less(i, j int) bool {
	return s.by(&s.entries[i], &s.entries[j])
}

type logEntrySorter struct {
	entries []LogEntry
	by      func(e1, e2 *LogEntry) bool
}

// RowKey to identify the row data in our table
type RowKey struct {
	id   int
	date int
}

// Take in a regular expression with text and return the matched text as a rectangle
func getRegexLogMatch(line string) (LogEntry, error) {

	r, err := regexp.Compile(".([0-9]+)-([0-9]+)-([0-9]+) ([0-9]+):([0-9]+). (.*)")
	check(err)

	allSubmatches := r.FindAllStringSubmatch(line, -1)
	ye, mo, da, ho, mi, log := "0", "0", "0", "0", "0", "LOG"
	ye, mo, da, ho, mi, log =
		allSubmatches[0][1], allSubmatches[0][2], allSubmatches[0][3], allSubmatches[0][4], allSubmatches[0][5], allSubmatches[0][6]

	year, err := strconv.Atoi(ye)
	month, err := strconv.Atoi(mo)
	day, err := strconv.Atoi(da)
	hour, err := strconv.Atoi(ho)
	minute, err := strconv.Atoi(mi)
	if err != nil {
		return LogEntry{}, err
	}

	date := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	return LogEntry{date, log}, nil

}

// Streamline error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Parse the input for the rectangle's distance from the left, distance from the top, width, and height
func parseLine(line string) (LogEntry, error) {

	entry, err := getRegexLogMatch(line)
	if err != nil {
		return LogEntry{}, err
	}
	return entry, nil

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
	allLogEntries := make([]LogEntry, 0)
	for fileScanner.Scan() {

		// We need to convert each line into a rectangle
		line := fileScanner.Text()
		rect, err := parseLine(line)
		if err != nil {
			continue
		}

		// Add the rectangle to the list
		allLogEntries = append(allLogEntries, rect)

	}

	// Sort the log entries by second, then minute
	date := func(e1, e2 *LogEntry) bool {
		return e1.date.Before(e2.date)
	}
	By(date).Sort(allLogEntries)

	for _, logEntry := range allLogEntries {
		fmt.Println(logEntry.date, ": ", logEntry.log)
	}

	// Output the answer
	end := time.Now()
	fmt.Println()
	fmt.Println("Time Taken: ", end.Sub(start))
	fmt.Println()

}
