#!/bin/bash

set -e

BIN=$1
WORK_DIR="$2/chains"
CHAIN_INDEX=$3
CHAIN_NODES=$4
PORT_START=$5
CHECK_BIN=$6
CHECK_CLEAN_BIN=$7
CASE_INDEX=$8

CHAIN_DIR="$WORK_DIR/chain_$CHAIN_INDEX"
MAX_NODE_INDEX=$(($CHAIN_NODES-1))
CHECK_PORT=$(($PORT_START+2000))
CHAIN_ID=60801
PRIV_KEYS=( 8d7d1f0bd3001d9ecbe81f3307807fbe117457fa6a85c936136b662e0f0d23c5 6c481be36abecba8acf9cfeed31831534fcf9589826f51777791ffda556f2ece 59e8ec1ce81ed2d37c35a51a157c5282ab6c27a33415d55979fa0002ff3da283 e5f78f05043ff7748f4e574f54c1fc53ba6749cf892fd6f3bc1fb47b19d4fb92 f85f544c4904b1c1bbacb64b10fbcc22ccf496f490d74e0f277017d00914eb54 906f247ab8735897c9b1d989bed46f85b6695049b9d330546281d25e25473225 9ad23fd74da1c61526b422410f0286f6948b8f8250cd514ec694615ecf5b0984 99227dabd235e2a6794a204206ebf01757414ca756d0ba3f61c44eb3269dfe67 12b86bd5b35363ec9217dbc176ebd54feceec60a1329099a2f310eaa1a5c512f 8259d475690e6916401da2c6f9297c48dfe1096a6957d81b12e69f22ed1591e6 848f52a89ae312698b29b286742e401d7e968e0495a57bb12a6a4af47344341a 6e446b5a867271ddd9c9a41ac031c019ba27cfc99310185981e7dd2ca1877b58)
PUB_KEYS=( 0x03e4294f6216a04e86a0608a4f6668f90d3bde6bd56f388821330f7bd217aa48f7 0x0211d4fa90e16c3cc69010b5ffee1f4bcd92e1f2a2adfbb52e91d0499f05874bfb 0x037bb0b29c54216eb8417c085af42a31d152349ec2716a61e701190f9c2aa1e0f3 0x036fcb8b3147d641f8c4e12aecdc800572e53350ef850abb4052972b5c6255f94f 0x03f2e32f9a060b8fe18736f5c4da328265d9d29ac13d5fed45649700a9c5f2cdca 0x026adcdfb010edff857a621342c29f264dd8ff4519826f68bb94542f655a844c4b 0x02620cbd4c10633f869e84b3e4f54e415d8bf8ead33e5d28db89d839e40b3627cc 0x0311c6f4cfa947b5e1beffc9c655eaf060676127e0e8bd6cf8e8722f10dc2598b5 0x033d0f77b6923e4e2cbb34e108c43f6e828acc01e5bcf0f85a4def11d1de871ad8 0x02a7be88bf49f46acabc091c7750faf6845fbdfdeab8f2a12f3264d49bc8a3a237 0x031aeb6af13ca1e70093eac2ff1928f2bd44dbdce0faa9b5d626869f03fb2a7f0a 0x023bf6b332c92ae44d6a330fa6f4d98040f077f1e7959a4f5a797263dc86235bd4)
ADDRESSES=(0x29F763Aa2acFC4688d7aB4aA07A9a78b7FFE8878 0xB6F16c5eb614e48a8AD99e8FAf4D1651528cB09C 0x1355D3A2f844a65c49c357E8c31a71D2E8C9084F 0x019Bd7499d9905942b105bdd0b144e003fdb3d68 0x40FBBE484b8Ee6139Af08446950B088e10b2306A 0x8C161d85fDC086AC6726bCEDe39f2CCB1Afa3bc8 0x22C21D4F64aabA7ec837b6B93639dB8cF514dAD5 0xBE6805C4c904B9cc39065ADC4cCcF4FCaB167AE6 0x250B0295816FB72F1e658AbCa07b0DDDCc0E7865 0xA406300B1F7BeA13158bb7DE117D6171e3B3396a 0xd2378898c87F72dE22eb67673fD23889796C694E 0x5C5bc95dD9F73E906BD2A036cc868Ddec6a266b7)

# clean db
$CHECK_CLEAN_BIN -case $CASE_INDEX

mkdir -p $WORK_DIR
rm -rf $CHAIN_DIR

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
          "0x011bD5423C5F77b5a0789E27f922535fd76B688F": {
              "balance": "300000000000000000000000"
            },
            "0x019Bd7499d9905942b105bdd0b144e003fdb3d68": {
              "balance": "300000000000000000000000"
            },
            "0x02484cb50AAC86Eae85610D6f4Bf026f30f6627D": {
              "balance": "300000000000000000000000"
            },
            "0x042a63149117602129B6922ecFe3111168C2C323": {
              "balance": "300000000000000000000000"
            },
            "0x0536896a5e38BbD59F3F369FF3682677965aBD19": {
              "balance": "300000000000000000000000"
            },
            "0x08135Da0A343E492FA2d4282F2AE34c6c5CC1BbE": {
              "balance": "300000000000000000000000"
            },
            "0x08A2DE6F3528319123b25935C92888B16db8913E": {
              "balance": "300000000000000000000000"
            },
            "0x09DB0a93B389bEF724429898f539AEB7ac2Dd55f": {
              "balance": "300000000000000000000000"
            },
            "0x0d803cdeEe5990f22C2a8DF10A695D2312dA26CC": {
              "balance": "300000000000000000000000"
            },
            "0x1003ff39d25F2Ab16dBCc18EcE05a9B6154f65F4": {
              "balance": "300000000000000000000000"
            },
            "0x11e8F3eA3C6FcF12EcfF2722d75CEFC539c51a1C": {
              "balance": "300000000000000000000000"
            },
            "0x1355D3A2f844a65c49c357E8c31a71D2E8C9084F": {
              "balance": "300000000000000000000000"
            },
            "0x145e2dc5C8238d1bE628F87076A37d4a26a78544": {
              "balance": "300000000000000000000000"
            },
            "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955": {
              "balance": "300000000000000000000000"
            },
            "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65": {
              "balance": "300000000000000000000000"
            },
            "0x16aDfbeFdEfD488C992086D472A4CA577a0e5e54": {
              "balance": "300000000000000000000000"
            },
            "0x1A7F6012DB5CD5f1065eA3fE58257765656796d8": {
              "balance": "50000000000000000000000"
            },
            "0x1BcB8e569EedAb4668e55145Cfeaf190902d3CF2": {
              "balance": "300000000000000000000000"
            },
            "0x1CBd3b2770909D4e10f157cABC84C7264073C9Ec": {
              "balance": "300000000000000000000000"
            },
            "0x1aac82773CB722166D7dA0d5b0FA35B0307dD99D": {
              "balance": "300000000000000000000000"
            },
            "0x1c87Bb9234aeC6aDc580EaE6C8B59558A4502220": {
              "balance": "300000000000000000000000"
            },
            "0x225356FF5d64889D7364Be2c990f93a66298Ee8D": {
              "balance": "300000000000000000000000"
            },
            "0x22C21D4F64aabA7ec837b6B93639dB8cF514dAD5": {
              "balance": "300000000000000000000000"
            },
            "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f": {
              "balance": "300000000000000000000000"
            },
            "0x250B0295816FB72F1e658AbCa07b0DDDCc0E7865": {
              "balance": "300000000000000000000000"
            },
            "0x2546BcD3c84621e976D8185a91A922aE77ECEc30": {
              "balance": "300000000000000000000000"
            },
            "0x25f825B63CA2a770A04712D7972A7d1989CFB196": {
              "balance": "50000000000000000000000"
            },
            "0x262595fa2a3A86adACDe208589614d483e3eF1C0": {
              "balance": "300000000000000000000000"
            },
            "0x265188114EB5d5536BC8654d8e9710FE72C28c4d": {
              "balance": "300000000000000000000000"
            },
            "0x29F763Aa2acFC4688d7aB4aA07A9a78b7FFE8878": {
              "balance": "300000000000000000000000"
            },
            "0x2c5b5Df7B681a1A13478811f2F8b6337fC8f7e1B": {
              "balance": "50000000000000000000000"
            },
            "0x2d34E0c18093Fff6c6aD119BdA9e9e54f3f1d4e1": {
              "balance": "50000000000000000000000"
            },
            "0x2f4f06d218E426344CFE1A83D53dAd806994D325": {
              "balance": "300000000000000000000000"
            },
            "0x30Bf53315437B47AeB9f6576F0f9094226342a58": {
              "balance": "41800000000000000000000000"
            },
            "0x344dca30F5c5f74F2f13Dc1d48Ad3A9069d13Ad9": {
              "balance": "300000000000000000000000"
            },
            "0x35304262b9E87C00c430149f28dD154995d01207": {
              "balance": "300000000000000000000000"
            },
            "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC": {
              "balance": "300000000000000000000000"
            },
            "0x3c3E2E178C69D4baD964568415a0f0c84fd6320A": {
              "balance": "300000000000000000000000"
            },
            "0x40FBBE484b8Ee6139Af08446950B088e10b2306A": {
              "balance": "300000000000000000000000"
            },
            "0x40Fc963A729c542424cD800349a7E4Ecc4896624": {
              "balance": "300000000000000000000000"
            },
            "0x4779d18931B35540F84b0cd0e9633855B84df7b8": {
              "balance": "300000000000000000000000"
            },
            "0x48fA7b63049A6F4E7316EB2D9c5BDdA8933BCA2f": {
              "balance": "300000000000000000000000"
            },
            "0x492C973C16E8aeC46f4d71716E91b05B245377C9": {
              "balance": "300000000000000000000000"
            },
            "0x49675089329c637c59b7DA465E00c7A8fe4c3247": {
              "balance": "2000000000000000000000000"
            },
            "0x4b23D303D9e3719D6CDf8d172Ea030F80509ea15": {
              "balance": "300000000000000000000000"
            },
            "0x4edFEDFf17ab9642F8464D6143900903dD21421a": {
              "balance": "300000000000000000000000"
            },
            "0x54ccCeB38251C29b628ef8B00b3cAB97e7cAc7D5": {
              "balance": "300000000000000000000000"
            },
            "0x553BC17A05702530097c3677091C5BB47a3a7931": {
              "balance": "300000000000000000000000"
            },
            "0x58621F440dA6cdf1A79C47E75f2dc3c804789A22": {
              "balance": "2000000000000000000000000"
            },
            "0x586BA39027A74e8D40E6626f89Ae97bA7f616644": {
              "balance": "300000000000000000000000"
            },
            "0x586aF62EAe7F447D14D25f53918814e04d3A5BA4": {
              "balance": "300000000000000000000000"
            },
            "0x5C5bc95dD9F73E906BD2A036cc868Ddec6a266b7": {
              "balance": "300000000000000000000000"
            },
            "0x5E661B79FE2D3F6cE70F5AAC07d8Cd9abb2743F1": {
              "balance": "300000000000000000000000"
            },
            "0x5F9805B99eCb2a9C3C23f9FDafF022efeeFD79b6": {
              "balance": "2000000000000000000000000"
            },
            "0x5eb15C0992734B5e77c888D713b4FC67b3D679A2": {
              "balance": "300000000000000000000000"
            },
            "0x61097BA76cD906d2ba4FD106E757f7Eb455fc295": {
              "balance": "300000000000000000000000"
            },
            "0x63c3686EF31C03a641e2Ea8993A91Ea351e5891a": {
              "balance": "300000000000000000000000"
            },
            "0x64492E25C30031EDAD55E57cEA599CDB1F06dad1": {
              "balance": "300000000000000000000000"
            },
            "0x6C8EA11559DFE79Ae3dBDD6A67b47F61b929398f": {
              "balance": "300000000000000000000000"
            },
            "0x6d69F301d1Da5C7818B5e61EECc745b30179C68b": {
              "balance": "300000000000000000000000"
            },
            "0x70997970C51812dc3A010C7d01b50e0d17dc79C8": {
              "balance": "300000000000000000000000"
            },
            "0x71F280DEA6FC5a03790941Ad72956f545FeB7a52": {
              "balance": "300000000000000000000000"
            },
            "0x71bE63f3384f5fb98995898A86B02Fb2426c5788": {
              "balance": "300000000000000000000000"
            },
            "0x73B3074ac649A8dc31c2C90a124469456301a30F": {
              "balance": "300000000000000000000000"
            },
            "0x7D86687F980A56b832e9378952B738b614A99dc6": {
              "balance": "300000000000000000000000"
            },
            "0x7Df8Efa6d6F1CB5C4f36315e0AcB82b02Ae8BA40": {
              "balance": "300000000000000000000000"
            },
            "0x7Ebb637fd68c523613bE51aad27C35C4DB199B9c": {
              "balance": "300000000000000000000000"
            },
            "0x7f8BB57E811B35783a125c5D14Fc2B820acA4C6B": {
              "balance": "2000000000000000000000000"
            },
            "0x81C0a104ebe7A0b20Aa039E13d54c97Fb43eb21a": {
              "balance": "50000000000000000000000"
            },
            "0x8263Fce86B1b78F95Ab4dae11907d8AF88f841e7": {
              "balance": "300000000000000000000000"
            },
            "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199": {
              "balance": "300000000000000000000000"
            },
            "0x86c53Eb85D0B7548fea5C4B4F82b4205C8f6Ac18": {
              "balance": "300000000000000000000000"
            },
            "0x87BdCE72c06C21cd96219BD8521bDF1F42C78b5e": {
              "balance": "300000000000000000000000"
            },
            "0x8C161d85fDC086AC6726bCEDe39f2CCB1Afa3bc8": {
              "balance": "300000000000000000000000"
            },
            "0x8C3229EC621644789d7F61FAa82c6d0E5F97d43D": {
              "balance": "300000000000000000000000"
            },
            "0x90F79bf6EB2c4f870365E785982E1f101E93b906": {
              "balance": "300000000000000000000000"
            },
            "0x924Ba5Ce9f91ddED37b4ebf8c0dc82A40202fc0A": {
              "balance": "300000000000000000000000"
            },
            "0x9394cb5f737Bd3aCea7dcE90CA48DBd42801EE5d": {
              "balance": "300000000000000000000000"
            },
            "0x976EA74026E726554dB657fA54763abd0C3a0aa9": {
              "balance": "300000000000000000000000"
            },
            "0x98D08079928FcCB30598c6C6382ABfd7dbFaA1cD": {
              "balance": "300000000000000000000000"
            },
            "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc": {
              "balance": "300000000000000000000000"
            },
            "0x9A50ed082Cf2fc003152580dcDB320B834fA379E": {
              "balance": "300000000000000000000000"
            },
            "0x9BF4beE5bfbEbb3a4b7060dAe40CA6fD49305D60": {
              "balance": "50000000000000000000000"
            },
            "0x9CDED4a330682cB88093214d1DE56D0B7aE525BF": {
              "balance": "2000000000000000000000000"
            },
            "0x9DCCe783B6464611f38631e6C851bf441907c710": {
              "balance": "300000000000000000000000"
            },
            "0x9E126C57330FA71556628e0aabd6B6B6783d99fA": {
              "balance": "300000000000000000000000"
            },
            "0x9eAF5590f2c84912A08de97FA28d0529361Deb9E": {
              "balance": "300000000000000000000000"
            },
            "0x9eF6c02FB2ECc446146E05F1fF687a788a8BF76d": {
              "balance": "300000000000000000000000"
            },
            "0xA1196426b41627ae75Ea7f7409E074BE97367da2": {
              "balance": "300000000000000000000000"
            },
            "0xA406300B1F7BeA13158bb7DE117D6171e3B3396a": {
              "balance": "300000000000000000000000"
            },
            "0xAb707cb80e7de7C75d815B1A653433F3EEc44c74": {
              "balance": "300000000000000000000000"
            },
            "0xB6F16c5eb614e48a8AD99e8FAf4D1651528cB09C": {
              "balance": "300000000000000000000000"
            },
            "0xBE6805C4c904B9cc39065ADC4cCcF4FCaB167AE6": {
              "balance": "300000000000000000000000"
            },
            "0xBc5BdceE96b1BC47822C74e6f64186fbA7d686be": {
              "balance": "300000000000000000000000"
            },
            "0xBcd4042DE499D14e55001CcbB24a551F3b954096": {
              "balance": "300000000000000000000000"
            },
            "0xC004e69C5C04A223463Ff32042dd36DabF63A25a": {
              "balance": "300000000000000000000000"
            },
            "0xC0543b0b980D8c834CBdF023b2d2A75b5f9D1909": {
              "balance": "300000000000000000000000"
            },
            "0xC39d4b38577b1d7a465A99FCd7Bd56ce1F821A5c": {
              "balance": "2000000000000000000000000"
            },
            "0xCcDd262f272Ee6C226266eEa13eE48D4d932Ce66": {
              "balance": "300000000000000000000000"
            },
            "0xD308a07F97db36C338e8FE2AfB09267781d00811": {
              "balance": "2000000000000000000000000"
            },
            "0xD3A01Ed4e1D0554a179Daebf95508b668767D441": {
              "balance": "2000000000000000000000000"
            },
            "0xD4A1E660C916855229e1712090CcfD8a424A2E33": {
              "balance": "300000000000000000000000"
            },
            "0xD6A098EbCc5f8Bd4e174D915C54486B077a34A51": {
              "balance": "300000000000000000000000"
            },
            "0xD9065f27e9b706E5F7628e067cC00B288dddbF19": {
              "balance": "300000000000000000000000"
            },
            "0xDFd99099Fa13541a64AEe9AAd61c0dbf3D32D492": {
              "balance": "300000000000000000000000"
            },
            "0xDe224dd66aeC53E5Ee234c94F9928601777dEac7": {
              "balance": "2000000000000000000000000"
            },
            "0xDf37F81dAAD2b0327A0A50003740e1C935C70913": {
              "balance": "300000000000000000000000"
            },
            "0xE5D3ab6883b7e8c35c04675F28BB992Ca1129ee4": {
              "balance": "300000000000000000000000"
            },
            "0xE74cEf90b6CF1a77FEfAd731713e6f53e575C183": {
              "balance": "300000000000000000000000"
            },
            "0xE77478D9E136D3643cFc6fef578Abf63F9Ab91B1": {
              "balance": "300000000000000000000000"
            },
            "0xEe7f6A930B29d7350498Af97f0F9672EaecbeeFf": {
              "balance": "300000000000000000000000"
            },
            "0xF0cE7BaB13C99bA0565f426508a7CD8f4C247E5a": {
              "balance": "300000000000000000000000"
            },
            "0xF0eeDDC5e015d4c459590E01Dcc2f2FD1d2baac7": {
              "balance": "300000000000000000000000"
            },
            "0xF23E054D8b4D0BECFa22DeEF5632F27f781f8bf5": {
              "balance": "300000000000000000000000"
            },
            "0xFABB0ac9d68B0B445fB7357272Ff202C5651694a": {
              "balance": "300000000000000000000000"
            },
            "0xFE0f143FcAD5B561b1eD2AC960278A2F23559Ef9": {
              "balance": "300000000000000000000000"
            },
            "0xa0EC9eE47802CeB56eb58ce80F3E41630B771b04": {
              "balance": "300000000000000000000000"
            },
            "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720": {
              "balance": "300000000000000000000000"
            },
            "0xa7580f28d5304b55594CfC1907F36D91b3D77cE5": {
              "balance": "2000000000000000000000000"
            },
            "0xa7e65F9aAb1b22ae918ec5e66Af94f2b1a35aB49": {
              "balance": "50000000000000000000000"
            },
            "0xaf69A66403Fa6CcF1ed658b75107BDE62126140d": {
              "balance": "50000000000000000000000"
            },
            "0xbDA5747bFD65F08deb54cb465eB87D40e51B197E": {
              "balance": "300000000000000000000000"
            },
            "0xbc8183bac3E969042736f7af07f76223D11D2148": {
              "balance": "300000000000000000000000"
            },
            "0xcBDc0F9a4C38f1e010bD3B6e43598A55D1868c23": {
              "balance": "300000000000000000000000"
            },
            "0xcC3a4d418e2062C977e4B642E3402520F77AeAc0": {
              "balance": "50000000000000000000000"
            },
            "0xcD9727b3605e879D2672a994F42F96EfEB5cB6ef": {
              "balance": "50000000000000000000000"
            },
            "0xcF2d5b3cBb4D7bF04e3F7bFa8e27081B52191f91": {
              "balance": "300000000000000000000000"
            },
            "0xcaBd7634D99020996c887AeE2B536f3fF1B71Fb6": {
              "balance": "2000000000000000000000000"
            },
            "0xcd3B766CCDd6AE721141F452C550Ca635964ce71": {
              "balance": "300000000000000000000000"
            },
            "0xd2378898c87F72dE22eb67673fD23889796C694E": {
              "balance": "300000000000000000000000"
            },
            "0xd98c495FE343dEDb29fB3Ed6b79834dA84f23631": {
              "balance": "2000000000000000000000000"
            },
            "0xdD2FD4581271e230360230F9337D5c0430Bf44C0": {
              "balance": "300000000000000000000000"
            },
            "0xdF3e18d64BC6A983f673Ab319CCaE4f1a57C7097": {
              "balance": "300000000000000000000000"
            },
            "0xdbD154FBDd5D174990c4b3F6d39E19aa31b1C5a4": {
              "balance": "50000000000000000000000"
            },
            "0xe141C82D99D85098e03E1a1cC1CdE676556fDdE0": {
              "balance": "300000000000000000000000"
            },
            "0xe8B1ff302A740fD2C6e76B620d45508dAEc2DDFf": {
              "balance": "300000000000000000000000"
            },
            "0xeed329483E4539a45Db863Df7b50436d50cC5276": {
              "balance": "50000000000000000000000"
            }
      },
    "governance": [
        {
          "Signer": "0x0Fdd8F55b123E4ce0E1DD3bB02dD9ac72c6adE7d",
          "Validator": "0x29F763Aa2acFC4688d7aB4aA07A9a78b7FFE8878"
        },
        {
          "Signer": "0xfEF9703F3486c31e75Aa191d15EE7347a55f3Da0",
          "Validator": "0xB6F16c5eb614e48a8AD99e8FAf4D1651528cB09C"
        },
        {
          "Signer": "0xBC1Ba91b59a05037A8f7bbc5Fc51B4f2a30B6275",
          "Validator": "0x1355D3A2f844a65c49c357E8c31a71D2E8C9084F"
        },
        {
          "Signer": "0xdCDf3AB480519b362931A552fcCE855b1Db07D79",
          "Validator": "0x019Bd7499d9905942b105bdd0b144e003fdb3d68"
        }
  ],
    "community_rate": 2000,
    "community_address": "0x79ad3ca3faa0F30f4A0A2839D2DaEb4Eb6B6820D",
    "coinbase": "0x0000000000000000000000000000000000000000",
    "difficulty": "0x1",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f89d801ef85494019bd7499d9905942b105bdd0b144e003fdb3d68941355d3a2f844a65c49c357e8c31a71d2e8c9084f9429f763aa2acfc4688d7ab4aa07a9a78b7ffe887894b6f16c5eb614e48a8ad99e8faf4d1651528cb09cb8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080",
    "gasLimit": "0x1fffffffffff",
    "nonce": "0x4510809143055965",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}
EOF

mkdir -p $NODE_DATA/geth

cat > "$NODE_DATA/static-nodes.json" << EOF
[
  "enode://e4294f6216a04e86a0608a4f6668f90d3bde6bd56f388821330f7bd217aa48f72b255f96eca30d4f7092da24f1d0b988066a93128c3abe8f9627a2d55bb30f1f@127.0.0.1:$(($PORT_START+1000))?discport=0",
  "enode://11d4fa90e16c3cc69010b5ffee1f4bcd92e1f2a2adfbb52e91d0499f05874bfbb5bf42bdadbc033d2b6115fd3d21752c3d41c64bebb0fb0871b2231015d16cbc@127.0.0.1:$(($PORT_START+1001))?discport=0",
  "enode://7bb0b29c54216eb8417c085af42a31d152349ec2716a61e701190f9c2aa1e0f3b5c61243a006aecb1fcd153aaf0c525d49e61fe709244bff453e89d023a1f31d@127.0.0.1:$(($PORT_START+1002))?discport=0",
  "enode://6fcb8b3147d641f8c4e12aecdc800572e53350ef850abb4052972b5c6255f94f018e7857fb1e04d692b9cadf1d7bacd150cbb8e1f5ea6e5fa9a7183cded368d5@127.0.0.1:$(($PORT_START+1003))?discport=0"
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

pushd $CHAIN_DIR

cat > check << EOF
#!/bin/bash
exec $CHECK_BIN -zion http://127.0.0.1:$PORT_START -port $CHECK_PORT -case $CASE_INDEX
EOF

chmod +x check

echo "cmd = ./check" > .wiz
echo "log=../case_check_$CASE_INDEX.log" >> .wiz
wizard start
popd

