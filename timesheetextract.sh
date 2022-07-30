#! /bin/sh

mycal=your.email@address.com

if [ -z "$1" ]; then
    echo "Provide the date or dates range as YYYY-MM-DD or YYYY-MM-DD_YYYY-MM-DD in a single argument"
else
    tsvf=$(icsreader $mycal $1 noclient true | grep exported | cut -f 4 -d ' ')
    converthourstobars $tsvf
    open $tsvf
fi
