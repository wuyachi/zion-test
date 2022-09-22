package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/polynetwork/bridge-common/util"
	"github.com/urfave/cli/v2"
)

const TEST_MM = "test test test test test test test test test test test junk"

type HDAddress struct {
	Index_1 uint32
	Index_2 uint32
}

func (hd *HDAddress) ToAddress() common.Address {
	privateKey := hd.PrivateKey()
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}
func (hd *HDAddress) PrivateKey() *ecdsa.PrivateKey {
	seed, err := hdwallet.NewSeedFromMnemonic(TEST_MM)
	if err != nil {
		panic(err)
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	m_44H, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		panic(err)
	}
	m_44H_60H, err := m_44H.Derive(hdkeychain.HardenedKeyStart + 60)
	if err != nil {
		panic(err)
	}
	m_44H_60H_I1H, err := m_44H_60H.Derive(hdkeychain.HardenedKeyStart + hd.Index_1)
	if err != nil {
		panic(err)
	}
	m_44H_60H_I1H_0, err := m_44H_60H_I1H.Derive(0)
	if err != nil {
		panic(err)
	}
	m_44H_60H_I1H_0_I2, err := m_44H_60H_I1H_0.Derive(hd.Index_2)
	if err != nil {
		panic(err)
	}
	privateKey, err := m_44H_60H_I1H_0_I2.ECPrivKey()
	if err != nil {
		panic(err)
	}
	return privateKey.ToECDSA()
}

type Discv5NodeID [64]byte

func (n Discv5NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

// PubkeyID returns a marshaled representation of the given public key.
func PubkeyID(pub *ecdsa.PublicKey) Discv5NodeID {
	var id Discv5NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}

func dump(ctx *cli.Context) (err error) {
	alloc := make(map[string]map[string]string)
	var validators []map[string]string
	var seeds, addresses []string
	unit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	allocPerUser := new(big.Int).Mul(big.NewInt(300000), unit)
	allocPerStaker := new(big.Int).Mul(big.NewInt(2000000), unit)
	allocPerValidator := new(big.Int).Mul(big.NewInt(300000), unit)
	allocPerProposal := new(big.Int).Mul(big.NewInt(50000), unit)
	left := new(big.Int).Mul(big.NewInt(100000000-300000*100-300000*12-50000*12-2000000*12), unit)

	var privs, pubs []string

	for i := 0; i < 100; i++ {
		w := &HDAddress{0, uint32(i + 1)}
		alloc[w.ToAddress().Hex()] = map[string]string{"balance": allocPerUser.String()}
	}
	for i := 0; i < 12; i++ {
		v := &HDAddress{uint32(i + 1), 1}
		alloc[v.ToAddress().Hex()] = map[string]string{"balance": allocPerValidator.String()}

		p := &HDAddress{uint32(i + 1), 3}
		alloc[p.ToAddress().Hex()] = map[string]string{"balance": allocPerProposal.String()}

		w := &HDAddress{uint32(i + 1), 4}
		alloc[w.ToAddress().Hex()] = map[string]string{"balance": allocPerStaker.String()}
		w.Index_2 = 1
		address := w.ToAddress().Hex()
		privs = append(privs, fmt.Sprintf("%x", crypto.FromECDSA(w.PrivateKey())))
		pubs = append(pubs, fmt.Sprintf("%#x", crypto.CompressPubkey(&w.PrivateKey().PublicKey)))
		addresses = append(addresses, address)

		if i > 7 {
			seeds = append(seeds, fmt.Sprintf("enode://%s@127.0.0.1:%v?discport=0", PubkeyID(&w.PrivateKey().PublicKey), 1000))
			w.Index_2 = 2
			validators = append(validators, map[string]string{
				"Validator": address,
				"Signer":    w.ToAddress().Hex(),
			})
		}
	}
	{
		w := &HDAddress{10000, 4}
		alloc[w.ToAddress().Hex()] = map[string]string{"balance": left.String()}
	}
	fmt.Println("allocation")
	fmt.Println(util.Verbose(alloc))
	fmt.Println("Validator private keys")
	fmt.Println(strings.Join(privs, " "))
	fmt.Println("Validator public keys")
	fmt.Println(strings.Join(pubs, " "))
	fmt.Println("Validator public addresses")
	fmt.Println(strings.Join(addresses, " "))
	fmt.Println("validators")
	fmt.Println(util.Verbose(validators))
	fmt.Println("Seeds")
	fmt.Println(util.Verbose(seeds))
	return
}
