#!/bin/bash
echo "Build urmmongo"
SRCPATH="$GOPATH/src/github.com/pascallimeux/urmmongo"
SRCBIN="/data/urmmongo/dist"

go build -ldflags "-s" $SRCPATH/server/urmmongo.go
if [ ! -d "$SRCBIN" ]; then
  mkdir $SRCBIN
fi
mv urmmongo $SRCBIN/urmmongo
cp $SRCPATH/server/config/config.json $SRCBIN/config.json
cp *.sh $SRCBIN
chmod u+x $SRCBIN/*.sh