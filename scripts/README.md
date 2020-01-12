# Scripts

This folder contains bash scripts automating the experiments.

## [clean.sh](clean.sh)

This script kills the running ipfs and ipfs-cluster-service processes. It is useful to run after a local simulation to kill the processes potentially slowing down the machine.

`Usage: ./clean.sh`

## [detergen.sh](detergen.sh)

This script generate a quasi-random topology to be deployed on deterlab. It creates the files [nodes.txt](../data/nodes.txt), with all individual nodes information, and [gen/](../data/gen) containing [details.txt](../data/gen/details.txt), the generation details record, and [deter.ns](../data/gen/deter.ns) the network topology to input on [deterlab](https://isi.deterlab.net).

```
Usage: ./detergen.sh [options]

Options:
  -D      maximum distance between nodes (default: 120)
  -K      number of ARA levels (default: 3)
  -N      number of nodes in the generated topology (default: 10)
  
Example: ./detergen.sh -D 150 -K 3 -N 15
```

## [prescript.sh](prescript.sh)

Cothority prescript copied on the remote hosts and executed at the begining of the program execution.
