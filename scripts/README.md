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

## [run.sh](run.sh)

This script runs the simulation. It reads the files created by [detergen.sh](detergen.sh) to get the details of the topology running on deterlab. It writes the script parameters to [details.txt](../data/details.txt), this file will be parsed by the simulation and contains the simulation parameters. Then it will run the experiment first in Vanilla mode and then in Cruxified mode, and output the results in the corresponding experiment folder in [results/](../results). Then, the simulation outputs are parsed to generate `pings.txt`, `min.txt`, `vanilla.txt` that are respectively the RTT between each pair of hosts, the interaction latency for cruxified IPFS and the interaction latency for vanilla IPFS.


```
./run.sh [options]
 
options:
  -h, --help                  show brief help
  -c, --cruxified             run only cruxified experiment
  -v, --vanilla               run only vanilla experiment
  -r, --remote                run simulation remotely
  -l, --local                 run simulation locally
  -p, --pings                 specify to compute new ping distances
  -m, --mode=MODE             specifiy ipfs-cluster mode (raft/crdt)
  -o, --operations=O          specify the number of operations to perform (int)
```

`-c` and `-v` flags specify that the experiment should be run respectively in Cruxified or Vanilla mode only. `-p` will make the simulation measure the ping distances again, otherwise it loads it from a previous experiment. `-m` can be either `raft` or `crdt` and specifies the consensus mode of IPFS Cluster. `-o` specifies the number of measurements.
