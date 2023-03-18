package main

import (
	"strings"
	"time"
)

var withZuluFormat string = "20060102T150405Z"
var noZuluFormat string = "20060102T150405"
var dateOnlyFormat string = "20060102"

func getLocalDateTimeFromString(dateString string) time.Time {
	// println("\n\n")
	// println("=============================")

	// test := "20221005T113000Z"
	// test := "20221005T113000"
	// dateString = "20221005"

	dateString = strings.TrimSpace(dateString)

	var parsed time.Time
	var err error

	if strings.HasSuffix(dateString, "Z") {
		parsed, err = time.Parse(withZuluFormat, dateString)
		// println("WITH ZULU: " + parsed.Local().String())
	} else if len(dateString) == 8 {
		parsed, err = time.ParseInLocation(dateOnlyFormat, dateString, time.Local)
		// println("DATE ONLY: " + parsed.Local().String())
	} else {
		parsed, err = time.ParseInLocation(noZuluFormat, dateString, time.Local)
		// println("NO ZULU: " + parsed.Local().String())
	}

	if err != nil {
		panic("unparsable datetime: " + dateString)
	}
	return parsed.Local()
}
