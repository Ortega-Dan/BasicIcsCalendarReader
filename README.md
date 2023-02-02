# BasicIcsCalendarReader

Extracts a .tsv (.csv like) report of hours (to open in Excel or other spreadsheet software) from an .ics calendar file.

When the events are prefixed with the `"CLIENT: Project, "` format, this generates the proper report of hours with columns for those classifications. 

To keep that report visually with some shortcuts and time control features from a Google Calendar, it is recommended to use [CalendarWorkHoursCounter](https://github.com/Ortega-Dan/CalendarWorkHoursCounter)

___
## Build:
### To native image [Recommended] (Requires GraalVM with Native Image feature):
```bash
./buildToNativeWithGraalVM.sh
```

### To regular jar:
```bash
javac BasicICSreader.java
jar -cvfe icsreader.jar BasicICSreader BasicICSreader.class
```


## Install: 
- If native-image built, just put the built `icsreader` in your path.
- If regular jar built, put the `icsreader.sh` in your path, remove the .sh extension (also put the `icsreader.jar` in the same directory preferably), and adjust the `icsreader` shell file in line 3 to point to the full path of the jar.


## Usage:
Run it without arguments to see usage.\
To simplify usage, the included `./timesheetextract.sh` can be added to the path (adjusted in line 5) and used as interface. (It also supports an alternative time-to-bars quick visual help about the length of the events. Which is found with instruction to build [here](convertHoursToBars/README.md))

___

# Prebuilt binaries:
Provided for both `icsreader` and `converthourstobars` ([ref](convertHoursToBars/README.md)) in releases for:
- Intel-Mac (amd64) (icsreader.jar instead of native built)
- M1-Mac (arm64)
- Linux (amd64)
- Windows (amd64) (even when documented scripts may need minor adjustments for Windows)