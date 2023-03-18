# BasicIcsCalendarReader

Extracts a .tsv (.csv like) report of hours (to open in Excel or other spreadsheet software) from an .ics calendar file.

When the events are prefixed with the `"CLIENT: Project, "` format, this generates the proper report of hours with columns for those classifications. 

To keep that report visually with some shortcuts and time control features from a Google Calendar, it is recommended to use [CalendarWorkHoursCounter](https://github.com/Ortega-Dan/CalendarWorkHoursCounter). You can export the .ics file from your Google Calendar from Google Calendar settings [like this](img/exportFromGcal.png).

___
Instructions are given for bash shell, but equivalent scripts are here with the .cmd file extension for Windows.

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
- For native-image built, just put the built `icsreader` in your path.
- If regular jar built, put the `icsreader.sh` in your path, remove the .sh extension for Linux or macOS (but do not remove the .cmd extension for Windows), also put the `icsreader.jar` in the same directory preferably, and adjust the `icsreader` script file in line 3 to point to the full path of the jar.


## Usage:
Run `ìcsreader` without arguments to see usage.\
To simplify usage, the included `./timesheetextract.sh` can be added to the path (adjusted in line 5) and used as interface. (It also supports an alternative time-to-bars quick visual help about the length of the events. Which is found with instruction to build [here](convertHoursToBars))

