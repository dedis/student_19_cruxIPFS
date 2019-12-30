#!/bin/sh

cat output.txt | grep details > details.txt
cat output.txt | grep minoptime > min.txt
cat output.txt | grep maxoptime > max.txt
cat output.txt | grep "ping node" > out_pings.txt

echo "Done!"
