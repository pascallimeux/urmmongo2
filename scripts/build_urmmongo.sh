#!/bin/bash
echo "Build urmmongo2"
. env.sh

go build -ldflags "-s" $SRCPATH/urmmongo.go

if [ ! -d "$SRCBIN" ]; then
  sudo mkdir -p $SRCBIN
  sudo chown -R $USER $DATAREPO
fi

mv urmmongo $SRCBIN/urmmongo
cp $SRCPATH/config.json $SRCBIN/config.json
cp *.sh $SRCBIN
chmod u+x $SRCBIN/*.sh

