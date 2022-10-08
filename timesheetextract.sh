#! /bin/sh
# custom runner helper

# replace with your common calendar ICS file
mycalfile=your.email@address.com.ics
# or use an env variable
# mycalfile=$mycal

if [ -z "$1" ]; then
    echo "Provide the date or dates range as YYYY-MM-DD or YYYY-MM-DD_YYYY-MM-DD in a single argument"
else
    # adjust with your prefered default usage
    tsvf=$(icsreader $mycalfile $1 noclient true | grep exported | cut -f 4 -d ' ')
    converthourstobars $tsvf
    # open with default application
    # on linux
    xdg-open $tsvf
    # on mac
    # open $tsvf
fi
