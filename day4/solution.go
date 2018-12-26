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

// guardLog to hold each guard's log data
type guardLog struct {
	id       string
	statuses []bool
	dates    []time.Time
	sleepSum int
}

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

// Parse the log file for an id or a status
func parseLog(log string) (string, bool) {

	rNum, err := regexp.Compile("([0-9]+)")
	check(err)

	numSubmatches := rNum.FindAllStringSubmatch(log, -1)
	if len(numSubmatches) >= 1 && len(numSubmatches[0]) > 1 {
		id := numSubmatches[0][1]
		return id, false
	}

	rBool, err := regexp.Compile("(falls)")
	check(err)

	boolSubmatches := rBool.FindAllStringSubmatch(log, -1)
	isAsleep := len(boolSubmatches) >= 1 && len(boolSubmatches[0]) > 1
	if isAsleep {
		return "", true
	}
	return "", false

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

	// Sort the log entries by date
	date := func(e1, e2 *LogEntry) bool {
		return e1.date.Before(e2.date)
	}
	By(date).Sort(allLogEntries)

	// Now that the log entries are sorted, we can parse the log data for times when guards were awake or asleep
	id := ""
	idToGuardLog := make(map[string]guardLog, 0)
	for _, logEntry := range allLogEntries {

		// We only get one piece of information at a time (id or status)
		parsedID, isAsleep := parseLog(logEntry.log)
		if parsedID != "" {
			id = parsedID
			if _, ok := idToGuardLog[id]; !ok {
				idToGuardLog[id] = guardLog{id, make([]bool, 0), make([]time.Time, 0), 0}
			}
			continue
		}

		// Update the statuses and dates on the guard log per entry
		log := idToGuardLog[id]
		idToGuardLog[id] = guardLog{id, append(log.statuses, isAsleep), append(log.dates, logEntry.date), 0}

	}

	// Keep track of the number of times a guard sleeps on any given minute
	idToMinuteToSleepCount := make(map[string]map[int]int, 0)
	for id, log := range idToGuardLog {

		// Create a map for the id if it doesn't exist
		if _, ok := idToMinuteToSleepCount[id]; !ok {
			idToMinuteToSleepCount[id] = make(map[int]int, 0)
		}

		// Each state has a corresponding time, which we can use to map out minutes guards are asleep
		for i, isAsleep := range log.statuses {

			if isAsleep {
				timeAsleep := log.dates[i+1].Sub(log.dates[i]).Minutes()
				for minute := 0; minute < int(timeAsleep); minute++ {

					// Update the id -> minute -> sleep count map
					date := log.dates[i].Add(time.Minute * time.Duration(minute))
					minute := date.Minute()
					idToMinuteToSleepCount[id][minute]++

					// Update the sum of time slept
					log := idToGuardLog[id]
					log.sleepSum++
					idToGuardLog[id] = log

				}
			}
			continue

		}

	}

	// Find the id with the highest number of minutes slept
	maxID := ""
	maxCount := 0
	for id, log := range idToGuardLog {
		if maxCount < log.sleepSum {
			maxCount = log.sleepSum
			maxID = id
		}
	}

	idInt, err := strconv.Atoi(maxID)
	check(err)

	// Find the highest count based on minute for the id
	maxCount = 0
	minute := 0
	for m, sleepCount := range idToMinuteToSleepCount[maxID] {
		if maxCount < sleepCount {
			maxCount = sleepCount
			minute = m
		}
	}

	// Output the answer
	end := time.Now()
	fmt.Println()
	fmt.Println("ID * Count: ", idInt*minute)
	fmt.Println("Time Taken: ", end.Sub(start))
	fmt.Println()

}
