# BasicIcsCalendarReader

Extracts a .tsv (.csv like) report of hours (to open in Excel or other spreadsheet software) from an .ics calendar file.

When the events are prefixed with "CLIENT: Project, " format, this generates the proper report of hours with columns for those classifications. 

To keep that report visually with some shortcuts and time control features from a Google Calendar, it is recommended to use [CalendarWorkHoursCounter](https://github.com/Ortega-Dan/CalendarWorkHoursCounter)

___
## Build:
### To native image (Requires GraalVM with Native Image feature):
```bash
./buildToNativeWithGraalVM.sh
```

### To regular jar:
```bash
javac BasicICSreader.java
jar -cvfe icsreader.jar BasicICSreader BasicICSreader.class
```


## Install: 
If native image built, just put it in your path.\
If regular jar build, put both `icsreader.sh` and `icsreader.jar` in your path.


## Usage:
Run it without arguments to see usage.\
Or to simplify usage the included `./timesheetextract.sh` can be added to the path adjusted and used as interface. (It also supports an alternative time-to-bars little visual help about the length of the events. and is found with instruction to build [here](convertHoursToBars/README.md))