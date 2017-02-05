#!/bin/bash
echo "make a tar.gz for urmmongo2 and send it to the integration server"
sh env.sh
tar cvzf urmmongo2.tar.gz $SRCBIN
scp urmmongo2.tar.gz orange@10.194.18.46:/tmp
rm urmmongo2.tar.gz