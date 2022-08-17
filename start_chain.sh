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
    NODE_ID="node_$i"
    P2P_PORT=$(($NODE_PORT+1000))
    logger -s start chain $i : $NODE_DIR ($P2P_PORT, $NODE_PORT)
    mkdir -p $NODE_DIR
    pushd $NODE_DIR

cat > node << EOF
$BIN --mine --miner.threads 1 \
--miner.etherbase=$COIN_BASE \
--identity=$NODE_ID \
--maxpeers=100 \
--syncmode full \
--gcmode archive \
--allow-insecure-unlock \
--datadir $NODE_DATA \
--networkid $CHAIN_ID \
--http.api admin,eth,debug,miner,net,txpool,personal,web3 \
--http --http.addr 127.0.0.1 --http.port $NODE_PORT --http.vhosts "*" \
--rpc.allow-unprotected-txs \
--nodiscover \
--port $P2P_PORT \
--verbosity 5
EOF
    chmod +x node


cat > "static-nodes.json" << EOF

EOF

cat > "genesis.json" << EOF
{
    "config": {
        "chainId": $CHAIN_ID, 
        "homesteadBlock": 0,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "byzantiumBlock": 0,
        "constantinopleBlock": 0,
        "petersburgBlock": 0,
        "istanbulBlock": 0,
        "berlinBlock": 0,
        "londonBlock": 0,
        "hotstuff": {
            "protocol": "basic"
        }
    },
    "alloc": {
        "0x5Ed9a6713962f04DA057e6A949394e002855DF72": {"balance": "30000000000000000000000000"},
        "0xb06Bf71465eA8071e83a87f6b914f8FFA6829f4b": {"balance": "10000000000000000000000000"},
        "0x328C78eb1b265F380381879e8F332fbB622EBAe5": {"balance": "10000000000000000000000000"},
        "0xfEb056933BE960a183a0178249f7195AEbB74A4C": {"balance": "10000000000000000000000000"},
        "0x51d352512e592f7c50880a1fE280259fD68B07a2": {"balance": "10000000000000000000000000"},
        "0xF6AE2D66B8a30104360a0188A7E8Da98eb336075": {"balance": "10000000000000000000000000"},
        "0x347Da040f1079C40AA76A84DeD08028D3425CC9D": {"balance": "10000000000000000000000000"},
        "0x3d7ADdF663a30Ef02E086aE9eC37c311C892AFF5": {"balance": "10000000000000000000000000"}
    },
    "governance": [
        {"Validator": "0x258af48e28e4a6846e931ddff8e1cdf8579821e5", "Signer": "0x5Ed9a6713962f04DA057e6A949394e002855DF72"},
        {"Validator": "0x6a708455c8777630aac9d1e7702d13f7a865b27c", "Signer": "0xb06Bf71465eA8071e83a87f6b914f8FFA6829f4b"},
        {"Validator": "0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae", "Signer": "0x328C78eb1b265F380381879e8F332fbB622EBAe5"},
        {"Validator": "0xad3bf5ed640cc72f37bd21d64a65c3c756e9c88c", "Signer": "0xfEb056933BE960a183a0178249f7195AEbB74A4C"}
    ],
    "community_rate": 2000,
    "community_address": "0x79ad3ca3faa0F30f4A0A2839D2DaEb4Eb6B6820D",
    "coinbase": "0x0000000000000000000000000000000000000000",
    "difficulty": "0x1",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f89c80f85494258af48e28e4a6846e931ddff8e1cdf8579821e5946a708455c8777630aac9d1e7702d13f7a865b27c948c09d936a1b408d6e0afaa537ba4e06c4504a0ae94ad3bf5ed640cc72f37bd21d64a65c3c756e9c88cb8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080",
    "gasLimit": "0x1fffff",
    "nonce": "0x4510809143055965",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}
EOF

cat > "nodekey.json" << EOF

EOF

    echo "cmd = ./node\nlog=../$i.log" > .wiz
    wizard start
    popd
done

