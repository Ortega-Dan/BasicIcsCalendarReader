#! /bin/sh
# custom runner helper

# replace with your common calendar ICS file
CAL_EMAIL=your.email@address.com.ics
# or remove previous line and use an env variable
# CAL_EMAIL=$MY_CAL_EMAIL_ENV_VAR
CAL_FILE="$CAL_EMAIL.ics"
CAL_ZIPFILE="$CAL_EMAIL.ical.zip"

if [ -z "$1" ]; then
    echo "Provide the date or dates range as YYYY-MM-DD or YYYY-MM-DD_YYYY-MM-DD in a single argument"
else
    unzip -o $CAL_ZIPFILE
    # adjust with your prefered default usage
    tsvf=$(icsreader $CAL_FILE $1 noclient false | grep "File exported to" | cut -f 4 -d ' ')
    withbars="$tsvf.withbars.tsv"
    converthourstobars $tsvf > $withbars
    # open with default application
    # on linux
    # xdg-open $withbars
    # on mac
    open $withbars
fi
