#!/bin/bash

set -e

WORK_DIR="$1/chains"
BIN=$2
CHAIN_INDEX=$3
CHAIN_NODES=$4
PORT_START=$5

MAX_NODE_INDEX=$(($CHAIN_NODES-1))
CHAIN_ID=100
COINBASE=0x5347440F0a74cA0FAfa95dD58eE8a15f6ABEd5c6


rm -rf $WORK_DIR

for i in $(seq 0 $MAX_NODE_INDEX) do
    NODE_DIR="$WORK_DIR/chain_$CHAIN_INDEX/node_$i"
    NODE_PORT=$((PORT_START+$i))
    NODE_DATA="data"
    P2P_PORT=$(($NODE_PORT+1000))
    logger -s stop chain $i : $NODE_DIR ($P2P_PORT, $NODE_PORT)
    mkdir -p $NODE_DIR
    pushd $NODE_DIR
    wizard stop
    popd
done