#!/bin/bash

f_K=0
f_N=0
f_D=0

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
            echo "-D                          specify the distance between nodes (default 120)"
            echo "-K                          specify the number of ARA levels (default 3)"
            echo "-N                          specify the number of nodes (default 10)"
            exit 0
            ;;

        -D)
            shift
            if test $# -gt 0; then
                f_D=$1
            else
                echo "no D specified"
                exit 1
            fi
            shift
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

        *)
            break
            ;;
    esac
done

if [[ $f_D -gt 0 ]]; then
    D=$f_D
else
    D=120
fi

if [[ $f_K -gt 0 ]]; then
    K=$f_K
else
    K=3
fi

if [[ $f_N -gt 0 ]]; then
    N=$f_N
else
    N=10
fi

cd ../detergen
go build
./detergen -N $N -R $N -K $K -SpaceMax $D
rm ../data/details.txt > /dev/null 2>&1
