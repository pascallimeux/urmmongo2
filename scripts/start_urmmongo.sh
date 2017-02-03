#!/bin/bash

USER="pascal"
URMMONGOPATH="/data/urmmongo/dist"
export URMMONGOCONFIGFILE="/data/urmmongo/dist/config.json"
LOGDIR="$URMMONGOPATH/logs"
export URMMONGOLOGFILE="$LOGDIR/urmmongo.log"

if [ ! -d "$LOGDIR" ]; then
	echo "Create log directory"
    sudo mkdir -p $LOGDIR
fi
    sudo chown -R $USER $LOGDIR

echo "Urmmongo process started."
CMD="$URMMONGOPATH/urmmongo"
eval "$CMD"
