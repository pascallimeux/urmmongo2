#!/bin/bash
echo "Stop urmmongo2"
PID=`pidof urmmongo`
if [ -n "$PID" ]
then
   kill -9 $PID
   echo "Urmmongo process stopped"
else
   echo "Could not send SIGTERM to kill urmmongo, probably it does not work:" >&2
fi