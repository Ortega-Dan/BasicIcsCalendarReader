package main

import (
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ortega-Dan/golang-stdin/stdin"
)

type pair struct {
	quarters int16
	text     string
}

func main() {

	if len(os.Args) == 1 {
		panic("Provide the path to a file as first argument")
	}

	lreader, err := stdin.NewLineByLineFileReader(os.Args[1])

	if err != nil {
		panic("Unable to open file")
	}

	var lines []pair

	var maxQuarters int16 = 0

	line, rerr := lreader()
	for rerr == nil {

		lineWithTask, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}.+", line)
		if err != nil {
			panic("Error in regex comparison")
		}

		if lineWithTask {
			splitted := strings.Split(line, "\t")
			hrs, err := strconv.ParseFloat(splitted[1], 64)
			if err != nil {
				panic("Error parsing hours")
			}

			q := int16(hrs / 0.25)
			t := strings.Join(splitted[1:], "\t")
			lines = append(lines, pair{q, t})

			if q > maxQuarters {
				maxQuarters = q
			}
		}

		line, rerr = lreader()
	}

	if rerr != io.EOF {
		panic("******* Error reading file *******")
	}

	for _, l := range lines {
		i := 0
		for ; i < int(l.quarters); i++ {
			print("|")
		}
		for ; i < int(maxQuarters); i++ {
			print(" ")
		}
		println("\t" + l.text)
	}
}
