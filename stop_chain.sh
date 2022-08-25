#!/bin/bash

set -e


BIN=$1
WORK_DIR="$2/chains"
CHAIN_INDEX=$3
CHAIN_NODES=$4
PORT_START=$5
CASE_INDEX=$6

CHAIN_DIR="$WORK_DIR/chain_$CHAIN_INDEX"
MAX_NODE_INDEX=$(($CHAIN_NODES-1))
CHAIN_ID=100



for i in $(seq 0 $MAX_NODE_INDEX)
do
    NODE_DIR="$WORK_DIR/chain_$CHAIN_INDEX/node_$i"
    NODE_PORT=$((PORT_START+$i))
    NODE_DATA="data"
    P2P_PORT=$(($NODE_PORT+1000))
    logger -s stop chain $i : $NODE_DIR \($P2P_PORT, $NODE_PORT\)
    mkdir -p $NODE_DIR
    pushd $NODE_DIR
    wizard stop
    popd
done

# rm -rf $WORK_DIR

pushd $CHAIN_DIR
wizard stop
popd
