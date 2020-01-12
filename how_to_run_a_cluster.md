# How to run an IPFS cluster

We suppose that IPFS and IPFS cluster (ipfs-cluster-service and ipfs-cluster-ctl) are already installed.

## Table of contents

- [How to run an IPFS cluster](#how-to-run-an-ipfs-cluster)
  - [Table of contents](#table-of-contents)
  - [Start IPFS nodes](#start-ipfs-nodes)
    - [Initialisation](#initialisation)
    - [Adapting `config`](#adapting-config)
    - [Starting ipfs daemon](#starting-ipfs-daemon)
    - [IPFS simple commands](#ipfs-simple-commands)
    - [How files are stored on IPFS](#how-files-are-stored-on-ipfs)
  - [Start IPFS Cluster](#start-ipfs-cluster)
    - [Cluster Initialisation](#cluster-initialisation)
    - [Adapting `service.json`](#adapting-servicejson)
    - [Starting a cluster](#starting-a-cluster)
  - [Cluster commands](#cluster-commands)
    - [Basic commands](#basic-commands)
  - [Scripts ideas](#scripts-ideas)

## Start IPFS nodes

### Initialisation

We need to start a distinct IPFS process for each node we want in the cluster. For each IPFS node, we need to specify the path to config files along with every command, in order to distinguish the different entities.

Initialising a node n:

```sh
>>> ipfs -c /path/to/config/ipfs_n init
initializing IPFS node at /path/to/config/ipfs_n
generating 2048-bit RSA keypair...done
peer identity: QmT4wQ6EeUsmEs8MngaZZXnP6wbNoPwbFBFziJyiXNeSpR
to get started, enter:

    ipfs cat /ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv/readme

```

This will generate the config files for the IPFS host.

```sh
>>> ls /path/to/config/ipfs_n/
api     config      datastore_spec  repo.lock
blocks  datastore   keystore        version
```

`api`  contains the api address of the IPFS host (default port is 5001).

```sh
>>> cat /path/to/config/ipfs_n/api
/ip4/127.0.0.1/tcp/5001
```

`config` is a json file that contains all the settings for the IPFS host, such as port numbers, peer id, secret key, max storage limit etc.

`blocks` folder contains the data blocks from the files that are stored on this IPFS instance. This folder is organised in subfolders of 2 characters (in base32), containing the data with name hash of the content in base32, i.e. CIQKU4GFO3TSSRGB4XYKCO3IZM2UFHF6MGNKEA47J57AC3ZUOPMY**AE**A.data will be stored in the folder named **AE**/.

`version` will obviously contain the version of the running IPFS.

The other files (`datastore_spec`, `repo.lock`) and folder (`datastore`, `keystore`) are here for some obscure reasons I did not figure out.

### Adapting `config`

For each IPFS node, we need to modify the port numbers that are used, otherwise all nodes will try to use the same ports. When created, all `config` files are exactly the same. We can do so either by directly writing in the config files or by executing config commands.

The ports used by IPFS are `4001` as swarm port, `5001` as API port and `8080` as gateway port. Those are the values to modify in order to be able to run multiple instances of IPFS on the same host.

The commands to update the ports are the following:

```sh
>>> ipfs -c /path/to/config/ipfs_n config --json Addresses.Swarm '["/ip4/0.0.0.0/tcp/4002","/ip6/::/tcp/4002"]'
>>> ipfs -c /path/to/config/ipfs_n config Addresses.API '/ip4/127.0.0.1/tcp/5002'
>>> ipfs -c /path/to/config/ipfs_n config Addresses.Gateway '/ip4/127.0.0.1/tcp/8081'
```

If the gateway address is *127.0.0.1* it is only accessible locally. If we want it to be accessible globally, we have to change the value to *0.0.0.0*.

We should also update the bootstraps ports, which are in a large json array which is hard to do manually. So we can simply use a script that replace the swarm port tcp/4001 everywhere with the new chosen port.

### Starting ipfs daemon

In order to start an ISPF instance, we simply have to use the following command:

```sh
>>> ipfs -c /path/to/config/ipfs_n daemon
```

### IPFS simple commands

Add, pin, cat, ls etc.

```sh
>>> ipfs -c /path/to/config/ipfs_n add /path/to/file.txt
added QmVBzZvMzrDeTT9wQzMkY7mciFckL6RenRRGHRCzPWR8SK file.txt
```

```sh
>>> ipfs -c /path/to/config/ipfs_n cat QmVBzZvMzrDeTT9wQzMkY7mciFckL6RenRRGHRCzPWR8SK
hello world!
```

### How files are stored on IPFS

Hash in base58 translated in base32 and put into folders. To translate the hash from base58 to base32, we use the [multibase-conv](https://github.com/multiformats/go-multibase) implementation in Go. The format corresponding to the given hash is `base58btc` (**z**) and we need to translate it to `base32` (**b**).

```sh
>>> multibase-conv b zQmVBzZvMzrDeTT9wQzMkY7mciFckL6RenRRGHRCzPWR8SK
bciqglsqjem5scgndtdx2r34aoj6zkyj2xhfmzgazamsykp32tyu5psa
```

Note: the `base32` names in IPFS system are capital letters, so we would need to capitalise the output string to get the perfect name match. The first character of the output (b) indicates the base (base32), so we should get rid of it to get the final hash `CIQGLSQJEM5SCGNDTDX2R34AOJ6ZKYJ2XHFMZGAZAMSYKP32TYU5PSA`.

The file containing the data of our example file will be stored at /path/to/config/ipfs_n/blocks/**PS**/CIQGLSQJEM5SCGNDTDX2R34AOJ6ZKYJ2XHFMZGAZAMSYKP32TYU5**PS**A.data. As mentionned before, the data file will be in the folder corresponding to the characters `(n-3)` and `(n-2)` with n being the length of the hash in `base32`.

## Start IPFS Cluster

### Cluster Initialisation

In order to create a local cluster, we need to run an ipfs-cluster-service instance for each ipfs node. To initialise an IPFS cluster instance, we need to run the following command:

```sh
>>> ipfs-cluster-service -c /path/to/ipfs-cluster/cluster_n init
20:03:57.202  INFO     config: Saving configuration config.go:347
ipfs-cluster-service configuration written to //path/to/ipfs-cluster/cluster_n/service.json
```

This command create the file `service.json` which is the equivalent of `config` for IPFS. All information concering the cluster instance, including ports are stored in this file.

### Adapting `service.json`

Similarly to the `config` file, `service.json` need to be adapted on every node to assign different ports. The ports that are used by ```ipfs-cluster-service``` are `5001` the same port as ipfs for API, `9094` for restAPI as http listen multiaddress, `9095` as ipfs proxy listen multaddress and `9096` as the main cluster multiaddress to communicate with other nodes.

As for `config`, the easiest way to update the ports is to run a simple script that changes the port numbers.

We also need to update the `secret` value to be the same everywhere. It is possible to joing a cluster only if the instance has the same `secret` value as the peers already in the cluster.

### Starting a cluster

To start a cluster, we should run the following command.

```sh
>>> ipfs-cluster-service -c /path/to/ipfs-cluster/cluster_n daemon

--- output ommitted ---

20:35:32.788  INFO    restapi: REST API (libp2p-http): ENABLED. Listening on:
        /p2p-circuit/ipfs/QmdWBarhotmeQ1VoFXFgBrYFwLjWVkkJvMNQESzZt4Fvve
        /ip4/127.0.0.1/tcp/9096/ipfs/QmdWBarhotmeQ1VoFXFgBrYFwLjWVkkJvMNQESzZt4Fvve
        /ip4/172.17.0.1/tcp/9096/ipfs/QmdWBarhotmeQ1VoFXFgBrYFwLjWVkkJvMNQESzZt4Fvve
        /ip4/172.16.0.254/tcp/9096/ipfs/QmdWBarhotmeQ1VoFXFgBrYFwLjWVkkJvMNQESzZt4Fvve
        /ip4/128.179.165.173/tcp/9096/ipfs/QmdWBarhotmeQ1VoFXFgBrYFwLjWVkkJvMNQESzZt4Fvve

 restapi.go:431
20:35:34.322  INFO  consensus: Current Raft Leader: QmdWBarhotmeQ1VoFXFgBrYFwLjWVkkJvMNQESzZt4Fvve raft.go:293
20:35:34.322  INFO    cluster: Cluster Peers (without including ourselves): cluster.go:406
20:35:34.322  INFO    cluster:     - No other peers cluster.go:408
20:35:34.322  INFO    cluster: ** IPFS Cluster is READY ** cluster.go:421
```

After the cluster is created, the file `peerstore` and the folder `raft/` will be created at `/path/to/ipfs-cluster/cluster_n/`. `peerstore` will contain the addresses of known peers when they join the cluster, and will be stored even after the new peers leave the cluster. `raft/` contains a file `raft.db` and an empty folder `snapshots`.

To join a cluster, we should run the same command as to start a cluster with the `--bootstrap` flag containing the address of a node in the cluster.

```sh
>>> ipfs-cluster-service -c /path/to/ipfs-cluster/cluster_n daemon --bootstrap /ip4/127.0.0.1/tcp/9096/ipfs/QmdWBarhotmeQ1VoFXFgBrYFwLjWVkkJvMNQESzZt4Fvve
```

After the node joined the cluster from the first time, it will store its peers addresses in the file `peerstore` and will be able to join the cluster without using the `--bootstrap` flag.

## Cluster commands

To run commands on the cluster once it is launched and multiple peers are connected, we have to use the command ```ipfs-cluster-ctl``` along with the node we want to run the command from. In a case of a local cluster, a node is identified by its restAPI port.

```sh
>>> ipfs-cluster-ctl --host "/ip4/127.0.0.1/tcp/9094" ...
```

### Basic commands

`id`, `peers`, `add`, `pin`, `status` etc.
