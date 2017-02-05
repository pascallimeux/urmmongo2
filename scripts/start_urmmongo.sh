#!/bin/bash
echo "Start urmmongo2"

# mandatory to systemctl service
. /data/urmmongo2/dist/env.sh

if [ ! -d "$LOGDIR" ]; then
	echo "Create log directory"
    sudo mkdir -p $LOGDIR
fi
    sudo chown -R $USER $LOGDIR

echo "Urmmongo2 process started."
CMD="$URMMONGOPATH/urmmongo"
eval "$CMD"
