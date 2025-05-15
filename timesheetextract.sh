#! /bin/sh
# custom runner helper

# replace with your common calendar ICS file
CAL_EMAIL=your.email@address.com.ics
CAL_PREFIX="" # a specific calender prefix, if needed (with the included trailing underscore _ to find it after unzipping). e.g. "Time Tracking_"

CAL_FILE="$CAL_PREFIX$CAL_EMAIL.ics"
CAL_ZIPFILE="$CAL_EMAIL.ical.zip"

if [ -z "$1" ]; then
    echo "Provide the date or dates range as YYYY-MM-DD or YYYY-MM-DD_YYYY-MM-DD in a single argument"
else
    unzip -o $CAL_ZIPFILE
    # adjust with your prefered default usage
    echo "Unzipped. Attempting to extract timesheet data from '$CAL_FILE'"
    tsvf=$(icsreader "$CAL_FILE" "$1" | grep "File exported to" | cut -f 4- -d ' ')
    echo "Extracted timesheet data to $tsvf"
    withbars="$tsvf.withbars.tsv"
    echo "Converting to bars format"
    converthourstobars "$tsvf" > "$withbars"
    # open with default application
    # on linux
    # xdg-open $withbars
    # on mac
    open "$withbars"
    open "$tsvf"
fi
