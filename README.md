# BasicIcsCalendarReader

(This is the supported and latest version written in **Golang** and previous unsupported version was moved to [legacy](legacy))
___

This command-line app extracts a .tsv (.csv like) report of hours (to open in Excel or other spreadsheet software) from an .ics calendar file.

When the events are prefixed with the `"CLIENT: Project, "` format, this generates the proper report of hours with columns for those classifications. 

To keep that report visually with some shortcuts and time control features from a Google Calendar, it is recommended to use [CalendarWorkHoursCounter](https://github.com/Ortega-Dan/CalendarWorkHoursCounter). You can export the .ics file from your Google Calendar from Google Calendar settings [like this](img/exportFromGcal.png).

___
Instructions are given for bash shell, but equivalent scripts are here with the .cmd file extension for Windows.

## To build or install from source:
Requires Golang 1.18+ installed locally:
```bash
# build with
go build
```
```bash
# or install with
go install
```


## Usage:
Run `ìcsreader` without arguments to see usage.\
To simplify usage, the included `./timesheetextract.sh` can be added to the path (adjusted in line 5) and used as interface. (It also supports an alternative time-to-bars quick visual help about the length of the events. Which is found with instruction to build [here](convertHoursToBars))

___

# Prebuilt binaries:
Provided for both `icsreader` and `converthourstobars` ([ref](convertHoursToBars)) in [releases](https://github.com/Ortega-Dan/BasicIcsCalendarReader/releases) for:
- Linux (amd64)
- M1-Mac (arm64)
- Intel-Mac (amd64)
- Windows (amd64)
