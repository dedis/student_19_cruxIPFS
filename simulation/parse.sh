#!/bin/sh


cat output.txt | grep optime > min.txt
cat output.txt | grep ping > tmp
tail -n +2 tmp > out_pings.txt
rm tmp

echo "Done!"
