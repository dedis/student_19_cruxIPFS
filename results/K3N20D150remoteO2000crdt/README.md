# K3N20D150remoteO2000crdt

This folder contains the results of the experiment run with Crux parameter `K=3`, `N=20` hosts, distance between hosts in the topology `D=150`, run on deterlab `remote`, with IPFS Cluster in `crdt` mode, and measuring `O=2000` pair interaction latencies.

## [data/](data)

This folder contains the data related with the experiment, such as the RTT between each pair of nodes and the interaction latency measurement for vanilla and cruxified IPFS.

## [graphs/](graphs)

This folder contains the graphs generated with the measurements from [data/](data).

## [details.txt](details.txt)

This file contains the details of the experiment.

## [ipfs.toml](ipfs.toml)

This file is the toml file that was used for the simulation.

## [nodes.txt](nodes.txt)

This file contains the nodes information, such as landmark level, X and Y coordinates and IP address.

## [output_c.txt](output_c.txt)

This file contains the output of the Cruxified IPFS simulation.

## [output_v.txt](output_v.txt)

This file contains the output of the Vanilla IPFS simulation.
