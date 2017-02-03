#!/bin/bash
sudo cp $GOPATH/src/github.com/pascallimeux/urmmongo/scripts/urmmongo.service /lib/systemd/system/
sudo chmod a+x /lib/systemd/system/urmmongo.service
sudo systemctl --system daemon-reload
sudo systemctl start urmmongo.service
sudo systemctl enable urmmongo.service
sudo systemctl is-active urmmongo.service
