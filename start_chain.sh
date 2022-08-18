#!/bin/bash

set -e

BIN=$1
WORK_DIR="$2/chains"
CHAIN_INDEX=$3
CHAIN_NODES=$4
PORT_START=$5

MAX_NODE_INDEX=$(($CHAIN_NODES-1))
CHAIN_ID=60801
PRIV_KEYS=(5f4873e69d20d714a9515dc3de5e343212a9714f5cafacfec828a0386373abc3 7b0a34d26b551ed9c30c4f9220cd3deb3e3047374fa9c447d7e1fd01795b2326 eea9fd97361a2d7361b968c24465eef3e563e77f0b56fbb2d05d3f13a3be073b 907b512c3eac860571a2299ebc495e6798d3fc650e559a7251959afff99bb8a3 dd7239d0f6ab62e998d1f89adb3173d1360ffe4d181e82019b51a379a0e9bf36 e11ae0f11e9d50cfe7aeb299f04238dd3f112c7851dbb05e35b2809e8a1f209a 8cff1dd500090e55f3834977f1917147e442c123c13c19c50a883fd11d6893b8 acadf1a084e579a5f80b07c0a8b7b2c8cdcddf9d0a50365df5f79c51f4f2c5e4 8c97a580029c0516df1f3dbf95b88143c9710f36602809c002b2489a7812b94e 2a4154119294f1648048358bcea3bbf0a2d98fd79248aaaff74e6fc9efdbdd74 f3bf5c79df5b8ada78cb0398c7df0a626d2e1c1307412e30e118a91529b34f7f c1234492837f05a2306ff9853125f3ca03df4dfe0fd76bf96ff25163421fd6cf)
PUB_KEYS=(0x03e60c67c8255b680975459484427a8dc039045db21a9684abaf999d8b20a16a1d 0x03742cdc052e86512d08a82424860e354cdc706716a2376834767351c2e5a8c18d 0x029a956ce28dbeb89afec1b0576be7e14daadbe636880668d57e504e0aa1e2dd43 0x03e213aad913ce29adc62df2a06123f69e7090c820c860f2fbc8fa376ed056e4d4 0x0285411291e1da075be73698a23043b7daf018fb54fb6070729e84d75213f218ac 0x03cd2c2a3251fd9170e39812e5291777f5cab17c9c777aebb50d686a2c402959b9 0x02a7ef7c7fc70aa13f461b04aed2d0a7d913f155687cabf44c8bd5251f2ff0d936 0x03461f090a1a5bf3fe4ec200bac3300864c99ec5cfdddacf8cf6a36939fb059353 0x032e758e402393ec28b21dbcca7eb464e20ac29733536145fc7beb88316df2f188 0x0334e2f8e04798e352e5b11a8b4c0a49a2fc21987bf74ecbd83291b9aa1852e367 0x0268317ebb2965c65f4471b1d215a6c5ce02692d10f271b8cac52439f69d992d21 0x0322c62a9ff1cc843fcccf2fc48da3cb3f928bab1cbf11348f9abefec345d7c063)
ADDRESSES=(0xa7580f28d5304b55594CfC1907F36D91b3D77cE5 0xD308a07F97db36C338e8FE2AfB09267781d00811 0x58621F440dA6cdf1A79C47E75f2dc3c804789A22 0xd98c495FE343dEDb29fB3Ed6b79834dA84f23631 0x49675089329c637c59b7DA465E00c7A8fe4c3247 0xDe224dd66aeC53E5Ee234c94F9928601777dEac7 0xD3A01Ed4e1D0554a179Daebf95508b668767D441 0xcaBd7634D99020996c887AeE2B536f3fF1B71Fb6 0x7f8BB57E811B35783a125c5D14Fc2B820acA4C6B 0xC39d4b38577b1d7a465A99FCd7Bd56ce1F821A5c 0x9CDED4a330682cB88093214d1DE56D0B7aE525BF 0x5F9805B99eCb2a9C3C23f9FDafF022efeeFD79b6)

rm -rf $WORK_DIR

for i in $(seq 0 $MAX_NODE_INDEX)
do
    NODE_DIR="$WORK_DIR/chain_$CHAIN_INDEX/node_$i"
    NODE_PORT=$((PORT_START+$i))
    NODE_DATA="data"
    NODE_ID="node_$i"
    P2P_PORT=$(($NODE_PORT+1000))
    logger -s start chain $i : $NODE_DIR \($P2P_PORT, $NODE_PORT\)
    mkdir -p $NODE_DIR
    pushd $NODE_DIR

cat > node << EOF
#!/bin/bash
exec $BIN --mine --miner.threads 1 \
--miner.etherbase=${ADDRESSES[$i]} \
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
        "0x02484cb50AAC86Eae85610D6f4Bf026f30f6627D": {
            "balance": "1000000000000000000000000"
        },
        "0x08135Da0A343E492FA2d4282F2AE34c6c5CC1BbE": {
            "balance": "1000000000000000000000000"
        },
        "0x08A2DE6F3528319123b25935C92888B16db8913E": {
            "balance": "1000000000000000000000000"
        },
        "0x09DB0a93B389bEF724429898f539AEB7ac2Dd55f": {
            "balance": "1000000000000000000000000"
        },
        "0x1003ff39d25F2Ab16dBCc18EcE05a9B6154f65F4": {
            "balance": "1000000000000000000000000"
        },
        "0x11e8F3eA3C6FcF12EcfF2722d75CEFC539c51a1C": {
            "balance": "1000000000000000000000000"
        },
        "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955": {
            "balance": "1000000000000000000000000"
        },
        "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65": {
            "balance": "1000000000000000000000000"
        },
        "0x1BcB8e569EedAb4668e55145Cfeaf190902d3CF2": {
            "balance": "1000000000000000000000000"
        },
        "0x1CBd3b2770909D4e10f157cABC84C7264073C9Ec": {
            "balance": "1000000000000000000000000"
        },
        "0x1aac82773CB722166D7dA0d5b0FA35B0307dD99D": {
            "balance": "1000000000000000000000000"
        },
        "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f": {
            "balance": "1000000000000000000000000"
        },
        "0x2546BcD3c84621e976D8185a91A922aE77ECEc30": {
            "balance": "1000000000000000000000000"
        },
        "0x2f4f06d218E426344CFE1A83D53dAd806994D325": {
            "balance": "1000000000000000000000000"
        },
        "0x30Bf53315437B47AeB9f6576F0f9094226342a58": {
            "balance": "26000000000000000000000000"
        },
        "0x35304262b9E87C00c430149f28dD154995d01207": {
            "balance": "1000000000000000000000000"
        },
        "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC": {
            "balance": "1000000000000000000000000"
        },
        "0x3c3E2E178C69D4baD964568415a0f0c84fd6320A": {
            "balance": "1000000000000000000000000"
        },
        "0x40Fc963A729c542424cD800349a7E4Ecc4896624": {
            "balance": "1000000000000000000000000"
        },
        "0x49675089329c637c59b7DA465E00c7A8fe4c3247": {
            "balance": "2000000000000000000000000"
        },
        "0x4b23D303D9e3719D6CDf8d172Ea030F80509ea15": {
            "balance": "1000000000000000000000000"
        },
        "0x553BC17A05702530097c3677091C5BB47a3a7931": {
            "balance": "1000000000000000000000000"
        },
        "0x58621F440dA6cdf1A79C47E75f2dc3c804789A22": {
            "balance": "2000000000000000000000000"
        },
        "0x5E661B79FE2D3F6cE70F5AAC07d8Cd9abb2743F1": {
            "balance": "1000000000000000000000000"
        },
        "0x5F9805B99eCb2a9C3C23f9FDafF022efeeFD79b6": {
            "balance": "2000000000000000000000000"
        },
        "0x5eb15C0992734B5e77c888D713b4FC67b3D679A2": {
            "balance": "1000000000000000000000000"
        },
        "0x61097BA76cD906d2ba4FD106E757f7Eb455fc295": {
            "balance": "1000000000000000000000000"
        },
        "0x70997970C51812dc3A010C7d01b50e0d17dc79C8": {
            "balance": "1000000000000000000000000"
        },
        "0x71bE63f3384f5fb98995898A86B02Fb2426c5788": {
            "balance": "1000000000000000000000000"
        },
        "0x7D86687F980A56b832e9378952B738b614A99dc6": {
            "balance": "1000000000000000000000000"
        },
        "0x7Ebb637fd68c523613bE51aad27C35C4DB199B9c": {
            "balance": "1000000000000000000000000"
        },
        "0x7f8BB57E811B35783a125c5D14Fc2B820acA4C6B": {
            "balance": "2000000000000000000000000"
        },
        "0x8263Fce86B1b78F95Ab4dae11907d8AF88f841e7": {
            "balance": "1000000000000000000000000"
        },
        "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199": {
            "balance": "1000000000000000000000000"
        },
        "0x86c53Eb85D0B7548fea5C4B4F82b4205C8f6Ac18": {
            "balance": "1000000000000000000000000"
        },
        "0x87BdCE72c06C21cd96219BD8521bDF1F42C78b5e": {
            "balance": "1000000000000000000000000"
        },
        "0x90F79bf6EB2c4f870365E785982E1f101E93b906": {
            "balance": "1000000000000000000000000"
        },
        "0x976EA74026E726554dB657fA54763abd0C3a0aa9": {
            "balance": "1000000000000000000000000"
        },
        "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc": {
            "balance": "1000000000000000000000000"
        },
        "0x9CDED4a330682cB88093214d1DE56D0B7aE525BF": {
            "balance": "2000000000000000000000000"
        },
        "0x9DCCe783B6464611f38631e6C851bf441907c710": {
            "balance": "1000000000000000000000000"
        },
        "0x9eAF5590f2c84912A08de97FA28d0529361Deb9E": {
            "balance": "1000000000000000000000000"
        },
        "0x9eF6c02FB2ECc446146E05F1fF687a788a8BF76d": {
            "balance": "1000000000000000000000000"
        },
        "0xBcd4042DE499D14e55001CcbB24a551F3b954096": {
            "balance": "1000000000000000000000000"
        },
        "0xC004e69C5C04A223463Ff32042dd36DabF63A25a": {
            "balance": "1000000000000000000000000"
        },
        "0xC39d4b38577b1d7a465A99FCd7Bd56ce1F821A5c": {
            "balance": "2000000000000000000000000"
        },
        "0xD308a07F97db36C338e8FE2AfB09267781d00811": {
            "balance": "2000000000000000000000000"
        },
        "0xD3A01Ed4e1D0554a179Daebf95508b668767D441": {
            "balance": "2000000000000000000000000"
        },
        "0xD4A1E660C916855229e1712090CcfD8a424A2E33": {
            "balance": "1000000000000000000000000"
        },
        "0xDe224dd66aeC53E5Ee234c94F9928601777dEac7": {
            "balance": "2000000000000000000000000"
        },
        "0xDf37F81dAAD2b0327A0A50003740e1C935C70913": {
            "balance": "1000000000000000000000000"
        },
        "0xFABB0ac9d68B0B445fB7357272Ff202C5651694a": {
            "balance": "1000000000000000000000000"
        },
        "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720": {
            "balance": "1000000000000000000000000"
        },
        "0xa7580f28d5304b55594CfC1907F36D91b3D77cE5": {
            "balance": "2000000000000000000000000"
        },
        "0xbDA5747bFD65F08deb54cb465eB87D40e51B197E": {
            "balance": "1000000000000000000000000"
        },
        "0xcF2d5b3cBb4D7bF04e3F7bFa8e27081B52191f91": {
            "balance": "1000000000000000000000000"
        },
        "0xcaBd7634D99020996c887AeE2B536f3fF1B71Fb6": {
            "balance": "2000000000000000000000000"
        },
        "0xcd3B766CCDd6AE721141F452C550Ca635964ce71": {
            "balance": "1000000000000000000000000"
        },
        "0xd98c495FE343dEDb29fB3Ed6b79834dA84f23631": {
            "balance": "2000000000000000000000000"
        },
        "0xdD2FD4581271e230360230F9337D5c0430Bf44C0": {
            "balance": "1000000000000000000000000"
        },
        "0xdF3e18d64BC6A983f673Ab319CCaE4f1a57C7097": {
            "balance": "1000000000000000000000000"
        },
        "0xe141C82D99D85098e03E1a1cC1CdE676556fDdE0": {
            "balance": "1000000000000000000000000"
        },
        "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266": {
            "balance": "1000000000000000000000000"
        }
    },
    "governance": [
        {
            "Signer": "0x2b382887D362cCae885a421C978c7e998D3c95a6",
            "Validator": "0xa7580f28d5304b55594CfC1907F36D91b3D77cE5"
        },
        {
            "Signer": "0xa5d32fADC0D92feBc80CfE80bB6992A23549A516",
            "Validator": "0xD308a07F97db36C338e8FE2AfB09267781d00811"
        },
        {
            "Signer": "0xAA81e0DC54fD1Bf47d1604EEa7C84c7A30d795c4",
            "Validator": "0x58621F440dA6cdf1A79C47E75f2dc3c804789A22"
        },
        {
            "Signer": "0xD21ea7573b76825dF8F0Ce82B4134F6333C693Ca",
            "Validator": "0xd98c495FE343dEDb29fB3Ed6b79834dA84f23631"
        }
    ],
    "community_rate": 2000,
    "community_address": "0x79ad3ca3faa0F30f4A0A2839D2DaEb4Eb6B6820D",
    "coinbase": "0x0000000000000000000000000000000000000000",
    "difficulty": "0x1",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f89d8014f8549458621f440da6cdf1a79c47e75f2dc3c804789a2294d308a07f97db36c338e8fe2afb09267781d0081194a7580f28d5304b55594cfc1907f36d91b3d77ce594d98c495fe343dedb29fb3ed6b79834da84f23631b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080",
    "gasLimit": "0x1fffff",
    "nonce": "0x4510809143055965",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}
EOF

mkdir -p $NODE_DATA/geth

cat > "$NODE_DATA/static-nodes.json" << EOF
[
  "enode://e60c67c8255b680975459484427a8dc039045db21a9684abaf999d8b20a16a1d3e78a73b268c146a5fefff657df416f31de36704aeca8c216ae7ef9a627a3e19@127.0.0.1:$(($PORT_START+1000))?discport=0",
  "enode://742cdc052e86512d08a82424860e354cdc706716a2376834767351c2e5a8c18d2d9a8b57bf0f54b151db9a5ff0158ff51832ae9e52d77e4f116cf30561eebb27@127.0.0.1:$(($PORT_START+1001))?discport=0",
  "enode://9a956ce28dbeb89afec1b0576be7e14daadbe636880668d57e504e0aa1e2dd43275831a98de73fcb13af099b288725fe3295196fdca6ffd032d79ca920a4c5ce@127.0.0.1:$(($PORT_START+1002))?discport=0",
  "enode://e213aad913ce29adc62df2a06123f69e7090c820c860f2fbc8fa376ed056e4d4c702694370ae7eee479e82210d387ec7d40f09ee78aeb1c7b03e1e38c44236ff@127.0.0.1:$(($PORT_START+1003))?discport=0"
]
EOF

echo ${PRIV_KEYS[$i]} > "$NODE_DATA/geth/nodekey"

logger -s init chain $CHAIN_INDEX
$BIN init genesis.json --datadir $NODE_DATA


    echo "cmd = ./node" > .wiz
    echo "log=../$i.log" >> .wiz
    wizard start
    popd
done

