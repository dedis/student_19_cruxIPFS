# results

This folder contains all the results of the conducted experiments on deterlab with the network in different settings. It also contains a python script to plot the measurements and a go program that gives interaction average latency for pairs of nodes according to the RTT between the nodes.

The title of the folders describe the content of the folder. I.e. the experience called `K3N15D100remoteO100raft` was conducted with Crux parameter `K=3`, `N=15` hosts, distance topology generation `D=100`, `remote` on deterlab, with `O=100` operations measured, in IPFS Cluster `raft` mode.

## [K3N20D150remoteO2000crdt/](K3N20D150remoteO2000crdt)

Results from this experiment is shown in the report, thus this folder is documented.

## [boxes.go](boxes.go)

Program showing the average interaction latency for pairs of nodes according to their RTT. The folder name, and boxes limit should be edited in the go program manually.

## [plot.py](plot.py)

Plot the graphs corresponding to the pair interaction latency according to the RTT between the writer and the reader node. The folder name should be edited manually in the python file.
