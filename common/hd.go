package common

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
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
