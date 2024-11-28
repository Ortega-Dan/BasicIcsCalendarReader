package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ortega-Dan/golang-stdin/stdin"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
)

/*
 * Reader for iCalendar 2.0 specification
 *
 * @author Dan Ortega
 */
func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Basic ICS Reader")
		fmt.Fprintln(os.Stderr, `Usage:
		
Dates must be in the format: YYYY-MM-DD
		
Usage arguments: [icsFilePath] [date or dateFrom_inclusiveDateTo] [clientsToFilter separated by pipes "|" (default: no filter)] [boolean isInclusiveFilter (default: true)]
		
		`)
		return
	}

	// ********** Setting up variables

	// Input file-path in .ics format
	inputFilePathString := os.Args[1]
	var inputFileBaseName string

	if !strings.EqualFold(filepath.Ext(inputFilePathString), ".ICS") {
		fmt.Fprintln(os.Stderr, "\nFirst argument (input file) must have an .ics extension")
		os.Exit(-1)
	} else {
		inputFileBaseName = strings.ReplaceAll(strings.ReplaceAll(filepath.ToSlash(filepath.Base(os.Args[1])), "/", ""), filepath.Ext(inputFilePathString), "")
	}

	// Dates from and to-inclusive (meaning from fromDate at 00:00 hours, until
	// toDate at 23:59 hours) in YYYY-MM-DD format
	datesMatcher := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})(_\d{4}-\d{2}-\d{2})?$`)
	dates := datesMatcher.FindString(os.Args[2])
	if dates == "" {
		fmt.Fprintln(os.Stderr, "\nError matching dates string. Make sure it is in the required format.\nExample: 2023-01-16 (or 2023-01-16_2023-01-31 for date ranges)")
		os.Exit(-1)
	}

	var stringDateFrom string
	var stringInclusiveDateTo string
	if strings.Contains(dates, "_") {
		datesArray := strings.Split(dates, "_")
		stringDateFrom = datesArray[0]
		stringInclusiveDateTo = datesArray[1]
	} else {
		stringDateFrom = dates
		stringInclusiveDateTo = dates
	}

	startRangeTime := getLocalDateTimeFromString(strings.ReplaceAll(stringDateFrom, "-", "") + "T000000")
	endRangeTime := getLocalDateTimeFromString(strings.ReplaceAll(stringInclusiveDateTo, "-", "") + "T235959")

	// ClientsToQuery: Case insensitive String of clients separated by pipes "|"
	// OR AN EMPTY STRING IF LOOKING FOR ALL ACTIVITIES in a time range or span.
	// For example to look for MTI and DATAFILE activities, set it to "mti|DataFile"
	clientsToQuery := ""

	// Set to false if you want to print all clients but the ones in the query ...
	// ... set to true if you want to print the clients in the query
	includeClientsInQuery := true

	if len(os.Args) > 3 {
		clientsToQuery = os.Args[3]
	}
	if len(os.Args) > 4 {
		includeClientsInQuery, _ = strconv.ParseBool(os.Args[4])
	}
	// ********** end of variables setting up
	//

	//
	// Logic goes from here ...
	clientFinderMatcher := regexp.MustCompile(`^(\w+\s*):.+`)
	projectFinderMatcher := regexp.MustCompile(`^(\w+\s*),.+`)
	clientsToQueryMatcher := regexp.MustCompile("(?i)(" + clientsToQuery + ").*")
	recordIgnoreShortcutMatcher := regexp.MustCompile(`\b(?i)(cal.ignore|c.ig)\b`)

	timeStartMatcher := regexp.MustCompile("DTSTART.*:")
	timeEndMatcher := regexp.MustCompile("DTEND.*:")

	// opening input file
	readLine, err := stdin.NewLineByLineFileReader(inputFilePathString)
	if err != nil {
		panic("Unable to open input file: " + inputFilePathString)
	}

	outFilePath := "." + string(os.PathSeparator) + regexp.MustCompile(`(.*/|.*\\|@.+)`).ReplaceAllString(inputFileBaseName, "") + "_" + dates + ".tsv"

	reportClientsToHours := treemap.NewWithStringComparator()
	otherClientsInRange := treeset.NewWithStringComparator()

	sortedReportEvents := treemap.NewWithStringComparator()

	var totalHours time.Duration
	var totalCount int64

	var line string
	var readingError error

	var versionFound bool
	for readingError == nil {

		startTime, endTime, summary, client, project := "", "", "", "", ""

		line, readingError = readLine()
		line = strings.TrimSpace(line)

		if !versionFound {
			if strings.HasPrefix(line, "VERSION:") {
				versionFound = true
				supportedVersionNumber := "2.0"
				fileVersionString := strings.TrimSpace(strings.ReplaceAll(line, "VERSION:", ""))

				fmt.Println("\niCalendar specification: " + fileVersionString + "\n")

				if supportedVersionNumber != fileVersionString {
					fmt.Fprintln(os.Stderr, "This program was written for iCalendar "+supportedVersionNumber)
					fmt.Fprintln(os.Stderr, "Please use an .ics file with "+supportedVersionNumber+" version or contact the developer of this program.")

					os.Exit(-1)
				}
			}
		}

		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			// var allEventLines strings.Builder

			line, readingError = readLine()
			line = strings.TrimSpace(line)
			for readingError == nil && !strings.HasPrefix(line, "END:VEVENT") {
				// allEventLines.WriteString(line)

				if strings.HasPrefix(line, "DTSTART") {
					startTime = strings.TrimSpace(timeStartMatcher.ReplaceAllString(line, ""))
				} else if strings.HasPrefix(line, "DTEND") {
					endTime = strings.TrimSpace(timeEndMatcher.ReplaceAllString(line, ""))
				} else if strings.HasPrefix(line, "SUMMARY:") {

					summary = line
					line, readingError = readLine()
					for readingError == nil && strings.HasPrefix(line, " ") {
						summary += line[1:]
						//----
						line, readingError = readLine()
					}

					summary = strings.TrimSpace(strings.ReplaceAll(strings.Replace(summary, "SUMMARY:", "", 1), `\,`, ","))

					// format event data

					clientMatches := clientFinderMatcher.FindStringSubmatch(summary)
					if clientMatches != nil {
						client = clientMatches[1]
						summary = strings.TrimSpace(strings.Replace(summary, client+":", "", 1))
						client = strings.TrimSpace(client)

						projectMatches := projectFinderMatcher.FindStringSubmatch(summary)
						if projectMatches != nil {
							project = projectMatches[1]
							summary = strings.TrimSpace(strings.Replace(summary, project+",", "", 1))
							project = strings.TrimSpace(project)
						}
					} else {
						client = "NoClient"
					}

				}

				// ----------------------
				if strings.HasPrefix(line, "END:VEVENT") {
					break
				}
				line, readingError = readLine()
			}

			// working with event data
			// Avoiding lunch or ignore-hours
			if strings.EqualFold(summary, "lunch") || recordIgnoreShortcutMatcher.MatchString(summary) {
				continue
			}

			if len(endTime) == 0 {
				continue
			}

			// Full days activities are only displayed and ignored after being filtered
			// see later usage of this variable
			isFullDayEvent := false
			if len(startTime) == 8 {
				startTime += "T000000"
				endTime += "T000000"
				isFullDayEvent = true
			}

			startDT := getLocalDateTimeFromString(startTime)
			endDT := getLocalDateTimeFromString(endTime)

			// Checking if date is in the range
			if (startDT.Equal(startRangeTime) || startDT.After(startRangeTime)) && (startDT.Before(endRangeTime) || startDT.Equal(endRangeTime)) {

				// Full days activities are only displayed and ignored after being filtered
				if isFullDayEvent {
					fmt.Println("FullDaysActivity: [" + client + "] " + startTime + "-" + endTime + " [" + summary + "]")
					continue
				}

				// client filtering (case insensitive matching)
				isClientMatching := clientsToQueryMatcher.MatchString(client)

				// if client matches and query is inclusive **** OR **** if client doesn't
				// match and query is exclusive
				if (isClientMatching && includeClientsInQuery) || (!isClientMatching && !includeClientsInQuery) {
					if startDT.Before(endDT) {
						duration := endDT.Sub(startDT)

						totalHours += duration
						totalCount += 1

						// Recording filtered clients to show them later
						previousDuration, found := reportClientsToHours.Get(client)
						if found {
							reportClientsToHours.Put(client, (previousDuration.(time.Duration))+duration)
						} else {
							reportClientsToHours.Put(client, duration)
						}

						// creating output line
						outLine := "\n" + strings.Split(startDT.String(), " ")[0] + "\t" + fmt.Sprintf("%.2f", duration.Hours()) + "\t" + client + "\t:" + project + ":\t" + summary

						sorterKey := startTime
						_, eventFound := sortedReportEvents.Get(sorterKey)
						for eventFound {
							sorterKey += "i"
							_, eventFound = sortedReportEvents.Get(sorterKey)
						}

						sortedReportEvents.Put(sorterKey, outLine)

					}
				} else {
					// Adding other clients in the dates-range for later information to the user
					otherClientsInRange.Add(client)
				}

			}
		}
	}
	if readingError != io.EOF {
		panic("******* Error reading file *******")
	}

	// concluding actions
	writeOutputFile(outFilePath, sortedReportEvents, totalHours)

	eventsWord := "events"
	if totalCount == 1 {
		eventsWord = "event"
	}
	eventsString := " (" + strconv.Itoa(int(totalCount)) + " " + eventsWord + ")"

	fmt.Println("\n *** Clients included in report: " + strconv.Itoa(reportClientsToHours.Size()) + eventsString + "\n")
	printSortedClientsAndHours(reportClientsToHours)

	fmt.Println("\nTotal: " + fmt.Sprintf("%.2f", totalHours.Hours()) + " hrs\n")

	if otherClientsInRange.Size() > 0 {

		fmt.Println("\n *** Other clients found in required dates range: " + strconv.Itoa(otherClientsInRange.Size()) + "\n")

		printSortedClients(otherClientsInRange)

	} else if reportClientsToHours.Size() == 0 {
		fmt.Print("\nNo clients found for filter and range.\n")
	} else {
		fmt.Print("\nAll found clients included in report.\n")
	}

	// the wording of this line is used by timesheetextract.sh
	fmt.Print("\n\nFile exported to: " + outFilePath + "\n\n")

}

func writeOutputFile(outFilePath string, sortedReportEvents *treemap.Map, totalHours time.Duration) {
	// creating output writer
	writer, err := os.Create(outFilePath)
	if err != nil {
		panic("Unable to create output file :" + outFilePath)
	}
	defer writer.Close()

	writer.WriteString("Date\tHours\tClient\tProject\tDescription")

	it := sortedReportEvents.Iterator()
	for it.Begin(); it.Next(); {
		writer.WriteString(it.Value().(string))
	}

	// Writing last line with total hours
	writer.WriteString("\n\t" + fmt.Sprintf("%.2f", totalHours.Hours()) + "\tTOTAL\t")
}

func printSortedClients(otherClientsInRange *treeset.Set) {

	it := otherClientsInRange.Iterator()
	for it.Begin(); it.Next(); {
		fmt.Println(it.Value())
	}
}

func printSortedClientsAndHours(reportClientsToHours *treemap.Map) {

	it := reportClientsToHours.Iterator()
	for it.Begin(); it.Next(); {
		fmt.Println(it.Key().(string) + ": " + fmt.Sprintf("%.2f", it.Value().(time.Duration).Hours()) + " hrs")
	}
}
