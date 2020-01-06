#!/bin/bash

usageMsg="usage: $0 filename vanilla/simulation"

case $# in
  [2] ) : nothing_may_be_OK ;;
  0 ) # no args, display usageMsg
      echo "$usageMsg" >&2 ; exit 1 ;;
esac

cd ..

mkdir ../../../reports/results/$1
mkdir ../../../reports/results/$1/data
mkdir ../../../reports/results/$1/graphs

n=$(cat $2/out_pings.txt | wc -l)
nodes=$(echo "sqrt($n)" | bc)
op=$(cat $2/min.txt | wc -l)

cat $2/details.txt > ../../../reports/results/$1/experiment.txt

echo "$nodes nodes" >> ../../../reports/results/$1/experiment.txt
echo "$op operations" >> ../../../reports/results/$1/experiment.txt

cp plot/*.pdf ../../../reports/results/$1/graphs
cp $2/out_pings.txt simulation/min.txt $2/max.txt ../../../reports/results/$1/data

rm plot/*.pdf
