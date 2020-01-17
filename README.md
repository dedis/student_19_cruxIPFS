# Crux IPFS
Author: [Guillaume Michel](https://github.com/guillaumemichel/)

Master Semester project: Speeding up the Inter-Planetary File System </br>
Laboratory: DEDIS, EPFL </br>
Project Advisor: Prof. Bryan Ford </br>
Project Supervisor: Cristina Basescu </br>
Student: Guillaume Michel

Presentation [slides](https://docs.google.com/presentation/d/1CBOIe3DP5Ju8UOPQKCX-kceS6qn9U05GWKWAMczlEbg/edit?usp=sharing) </br>
Project [report](report/report.pdf) </br>
IPFS and IPFS Cluster configuration [guide](how_to_run_a_cluster.md) </br>

## What is this repository

This repository contains Cruxified IPFS, an implementation of Crux [1] applied on [IPFS Cluster](https://cluster.ipfs.io/). Crux provides an upper bound to the interaction between any two nodes running IPFS, that is a small multiple of their RTT. This repository contains the interaction latency tests, that were performed on the platform [Deterlab](https://www.isi.deterlab.net/). I invite you to read the [report](report/report.pdf) for a better understanding.

## How to use this repository

First, make sure to have Go installed.

In order to generate a network topology, go to [scripts/](scripts) and run the script [detergen.sh](scripts/detergen.sh). Then upload the output file [deter.ns](data/gen/deter.ns) on the platform deterlab to create the simulated network, and swap in the experiment. Once the experiment is swapped in, you can run the script [run.sh](scripts/run.sh) to start the measurement on the experiment. If you run the experiment locally, you shouldn't need to have installed [ipfs](https://github.com/ipfs/go-ipfs/) nor [ipfs-cluster](https://github.com/ipfs/ipfs-cluster/), the prescript should install it on your machine. Feel free to contact me if you experience any issue.

## What is in this repository

### [data/](data)

Contains the data related to the current simulation.

### [detergen/](detergen)

Contains the code that generate a network topology to be deployed on deterlab.

### [gentree/](gentree)

Contains the code that generates ARAs, and create Onet trees for these.

### [operations/](operations)

Contains Cruxified IPFS client and performance tests.

### [reports/](reports)

Contains the source code and pdf of the report of the project.

### [results/](results)

Contains results of various simulation as well as tools to anaylze the results.

### [scripts/](scripts)

Contains the scripts used to run the simulation.

### [service/](service)

Contains the code used to interact with Onet, the Services, Protocols and IPFS configuration.

### [simulation/](simulation)

Contains the skeleton of the program.

## References

[1] Cristina Basescu, Michael F. Nowlan, Kirill Nikitin, Jose M. Faleiro, and Bryan Ford. Crux: Locality-Preserving Distributed Services. May 2018. arXiv, cs.DC:1405.0637v2. https://arxiv.org/pdf/1405.0637.pdf. </br>

Contact: guillaume.michel@epfl.ch
