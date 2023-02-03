@echo off
REM custom runner helper

REM replace with your common calendar ICS file
set mycalfile=ortega.dan2010@gmail.com.ics
REM or use an env variable
REM mycalfile=$mycal

IF [%1]==[] (
echo Provide the date or dates range as YYYY-MM-DD or YYYY-MM-DD_YYYY-MM-DD in a single argument
) ELSE (
icsreader %mycalfile% %1 noclient true

REM NOTE: this doesn't support triggering converthourstobars yet
REM you may run it manually after installing with: converthourstobars the-tsv-output-file.tsv

)