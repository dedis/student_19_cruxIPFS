#!/bin/bash

# initialize default values

d_K=3
d_N=10
d_D=120
d_ops=100
d_t=20

# initialize empty variables

f_local=false
f_remote=false
f_mode=''
f_ops=-1
f_t=0
f_pings=false
f_van=false
f_crux=false
f_N=0
f_K=0

# parsing parameters
while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "$0 - running the experiment"
      echo " "
      echo "$0 [options]"
      echo " "
      echo "options:"
      echo "-h, --help                  show brief help"
      echo "-c, --cruxified             run only cruxified experiment"
      echo "-v, --vanilla               run only vanilla experiment"
      echo "-r, --remote                run simulation remotely"
      echo "-l, --local                 run simulation locally"
      #echo "-K                          specify the number of ARA levels"
      #echo "-N                          specify the number of nodes"
      echo "-p, --pings                 specify to compute new ping distances"
      echo "-m, --mode=MODE             specifiy ipfs-cluster mode (raft/crdt)"
      echo "-o, --operations=O          specify the number of operations to perform (int)"
      exit 0
      ;;

    -K)
        shift
        if test $# -gt 0; then
            f_K=$1
        else
            echo "no K specified"
            exit 1
        fi
        shift
        ;;


    -l|--local)
      if $f_remote; then
        echo "Simulation cannot be local and remote at the same time ..."
        exit 1
      fi

      f_local=true
      shift
      ;;

    -m)
      shift
      if test $# -gt 0; then
        f_mode=$1
      else
        echo "no mode specified"
        exit 1
      fi
      shift
      ;;
    --mode*)
      f_mode=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;

    -N)
        shift
        if test $# -gt 0; then
            f_N=$1
        else
            echo "no N specified"
            exit 1
        fi
        shift
        ;;

    -o)
        shift
            if test $# -gt 0; then
                f_ops=$1
            else
                echo "No operation number specified"
                exit 1
            fi
        shift
        ;;
    --operations*)
        f_ops=`echo $1 | sed -e 's/^[^=]*=//g'`
        shift
        ;;

    -r|--remote)
      if $f_local; then
          echo "Simulation cannot be local and remote at the same time ..."
          exit 1
      fi
      f_remote=true
      shift
      ;;

    -p|--pings)
        f_pings=true
        shift
        ;;

    -t)
        shift
        if test $# -gt 0; then
            f_t=$1
        else
            echo "no timeout specified"
            exit 1
        fi
        shift
        ;;
    --timeout)
        f_t=`echo $1 | sed -e 's/^[^=]*=//g'`
        shift
        ;;

    -v|--vanilla)
        f_van=true
        shift
        ;;

    *)
      break
      ;;
  esac
done

rm ../data/output_c.txt > /dev/null 2>&1
rm ../data/output_v.txt > /dev/null 2>&1
rm ../data/results/* > /dev/null 2>&1

# exporting execution details to text file
DETFILE="../data/details.txt"
if ! [ -f "$DETFILE" ]; then

    GENFILE='../data/gen/details.txt'
    if ! [ -f "$GENFILE" ]; then
        # gen details do not exist
        if ! $f_remote; then

            if [[ $f_N -gt 0 ]]; then
                d_N=$f_N
            fi
            if [[ $f_K -gt 0 ]]; then
                d_K=$f_K
            fi

            echo "Running detergen with default parameters, local simulation"
            echo './detergen -N '$d_N' -R '$d_N' -K '$d_K' -SpaceMax '$d_D
            cd ../detergen
            go build
            ./detergen -N $d_N -R $d_N -K $d_K -SpaceMax $d_D
            cd ../scripts
        else
            echo $GENFILE' missing. Network topology should be generated and swapped in on deterlab before deploying remotely'
            exit 1
        fi
    fi

    # dealing with default parameters

    # local flag
    if $f_local; then
        remote=false
    else
        remote=true
    fi

    # operations number
    if [[ $f_ops -eq -1 ]]; then
        ops=$d_ops
    fi

    # mode: raft/crdt
    if [ "$f_mode" = "raft" ]; then
        mode=raft
    else
        mode=crdt
    fi

    # timeout
    if [[ $f_t -gt 0 ]]; then
        t=$f_t
    else
        t=$d_t
    fi

    pings=true

    # read detergen details
    K=`cat $GENFILE | grep K | cut -d '=' -f2`
    N=`cat $GENFILE | grep N | cut -d '=' -f2`
    R=`cat $GENFILE | grep R | cut -d '=' -f2`
    D=`cat $GENFILE | grep D | cut -d '=' -f2`

else
    # read detergen details
    K=`cat $DETFILE | grep K | cut -d '=' -f2`
    N=`cat $DETFILE | grep N | cut -d '=' -f2`
    R=`cat $DETFILE | grep R | cut -d '=' -f2`
    D=`cat $DETFILE | grep D | cut -d '=' -f2`
    remote=`cat $DETFILE | grep remote | cut -d '=' -f2`
    ops=`cat $DETFILE | grep ops | cut -d '=' -f2`
    mode=`cat $DETFILE | grep mode | cut -d '=' -f2`
    t=`cat $DETFILE | grep timeout | cut -d '=' -f2`

    # update variables according to flags
    if $f_local; then
        remote=false
    elif $f_remote; then
        remote=true
    fi

    if [[ $f_ops -ne -1 ]]; then
        ops=$f_ops
    fi

    if ! [ "$f_mode" = '' ]; then
        if [ "$f_mode" = 'raft' ] || [ "$f_mode" = 'crdt' ];then
            mode=$f_mode
        fi
    fi

    if [[ $f_N -gt 0 ]]; then
        N=$f_N
        R=$f_N
    fi
    if [[ $f_K -gt 0 ]]; then
        K=$f_K
    fi

    pings=$f_pings
fi

# generate ipfs.toml
printf 'Simulation = "IPFS"\nServers = '$N'\nBf = '$(($N-1))'\nRounds = 1\nSuite = "Ed25519"\nPrescript = "prescript.sh"\n\nDepth\n1' > ../simulation/ipfs.toml
printf 'Simulation = "IPFS"\nServers = '$N'\nBf = '$(($N-1))'\nRounds = 1\nSuite = "Ed25519"\nPrescript = "prescript.sh"\n\nDepth\n1' > ../data/ipfs.toml

cd ../simulation
go build

# run cruxified experiment
if ! $f_van; then

    cruxified=true

    # writing parameters to data/details.txt
    echo 'K='$K > $DETFILE
    echo 'N='$N >> $DETFILE
    echo 'R='$R >> $DETFILE
    echo 'D='$D >> $DETFILE
    echo 'remote='$remote >> $DETFILE
    echo 'ops='$ops >> $DETFILE
    echo 'pings='$pings >> $DETFILE
    echo 'mode='$mode >> $DETFILE
    echo 'cruxified='$cruxified >> $DETFILE

    output_c='../data/output_c.txt'

    echo 'Starting cruxified experiment'
    echo `cat $DETFILE`
    echo

    # run cruxified experiment
    if $remote; then
        ./simulation -platform deterlab -mport 10008 ipfs.toml > $output_c
    else
        ./simulation ipfs.toml > $output_c
    fi

    # wait for the simulation to finish
    while ! grep -q "Done" "$output_c"; do
        sleep 15
    done

    # parse output
    cat $output_c | grep "ping node" > ../data/results/pings.txt
    cat $output_c | grep minoptime > ../data/results/min.txt
    cat $output_c | grep maxoptime > ../data/results/max.txt

    pings=false
fi

# run vanilla experiment
if ! $f_crux; then

    cruxified=false
    # writing parameters to data/details.txt
    echo 'K='$K > $DETFILE
    echo 'N='$N >> $DETFILE
    echo 'R='$R >> $DETFILE
    echo 'D='$D >> $DETFILE
    echo 'remote='$remote >> $DETFILE
    echo 'ops='$ops >> $DETFILE
    echo 'pings='$pings >> $DETFILE
    echo 'mode='$mode >> $DETFILE
    echo 'cruxified='$cruxified >> $DETFILE

    output_v='../data/output_v.txt'

    echo 'Starting vanilla experiment'
    echo `cat $DETFILE`
    echo

    # run experiment
    if $remote; then
        ./simulation -platform deterlab -mport 10008 ipfs.toml > $output_v
    else
        ./simulation ipfs.toml > $output_v
    fi

    # wait for the simulation to finish
    while ! grep -q "Done" "$output_v"; do
        sleep 15
    done

    # parse output
    if $f_van; then
        cat $output_v | grep "ping node" > ../data/results/pings.txt
    fi
    cat $output_v | grep minoptime > ../data/results/vanilla.txt

fi

# plot graph
cd ../plot
rm *.pdf
python3 plot.py

# save results

mkdir '../../results' > /dev/null 2>&1

if $remote; then
    deploy=remote
else
    deploy=local
fi

path='../../results/K'$K'N'$N'D'$D$deploy'O'$ops$mode

mkdir $path
mkdir $path'/data'
mkdir $path'/graphs'

cp '../data/ipfs.toml' $path
cp '../data/nodes.txt' $path
cp '../data/details.txt' $path

cp -r '../plot/.' $path'/graphs'
rm $path'/graphs/plot.py'
cp -r '../data/results/.' $path'/data'

echo 'Results saved to '$path
