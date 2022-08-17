package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

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

func dump(ctx *cli.Context) (err error) {
	alloc := make(map[string]map[string]string)
	var validators []map[string]string
	unit := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	allocPerUser := new(big.Int).Mul(big.NewInt(1000000), unit)
	allocPerStaker := new(big.Int).Mul(big.NewInt(2000000), unit)
	left := new(big.Int).Mul(big.NewInt(100000000-1000000*20-2000000*12), unit)

	var privs, pubs []string

	for i := 0; i < 20; i++ {
		w := &HDAddress{0, uint32(i)}
		alloc[w.ToAddress().Hex()] = map[string]string{"balance": allocPerUser.String()}
	}
	for i := 0; i < 12; i++ {
		w := &HDAddress{uint32(i + 1), 4}
		address := w.ToAddress().Hex()
		alloc[address] = map[string]string{"balance": allocPerStaker.String()}
		privs = append(privs, fmt.Sprintf("%x", crypto.FromECDSA(w.PrivateKey())))
		pubs = append(pubs, fmt.Sprintf("%#x", crypto.FromECDSAPub(&w.PrivateKey().PublicKey)))

		if i < 4 {
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
	fmt.Println(util.Verbose(privs))
	fmt.Println("Validator public keys")
	fmt.Println(util.Verbose(pubs))
	fmt.Println("validators")
	fmt.Println(util.Verbose(validators))

	return
}
