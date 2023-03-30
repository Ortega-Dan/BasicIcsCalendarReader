#! /bin/sh
# custom runner helper

# replace with your common calendar ICS file
MYCALFILE=your.email@address.com.ics
# or remove previous line and use an env variable
# MYCALFILE=$mycal

if [ -z "$1" ]; then
    echo "Provide the date or dates range as YYYY-MM-DD or YYYY-MM-DD_YYYY-MM-DD in a single argument"
else
    # adjust with your prefered default usage
    tsvf=$(icsreader $MYCALFILE $1 noclient true | grep "File exported to" | cut -f 4 -d ' ')
    converthourstobars $tsvf
    # open with default application
    # on linux
    xdg-open $tsvf
    # on mac
    # open $tsvf
fi
