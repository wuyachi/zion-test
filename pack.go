package main

import (
	"math/big"

	"github.com/devfans/zion-sdk/contracts/native/utils"
	"github.com/polynetwork/bridge-common/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var DEFAULT_GAS_PRICE = big.NewInt(1000000000)
var DEFAULT_GAS_LIMIT uint64 = 10000000
var NODE_MANAGER_CONTRACT = utils.NodeManagerContractAddress
var ZION_CHAINID = big.NewInt(60801)

type MethodType int

const (
	QUERY MethodType = iota
	TX
	CHECK_BALANCE
)

type Param interface {
	Encode() ([]byte, error)
}

type RawCase struct {
	Index   int64
	Actions []*RawAction
}

type RawAction struct {
	Row              []string
	MethodName       string
	Input            Param
	ShouldSucceed    bool
	Assertions       []Assertion
	Sender           HDAddress
	CheckBalancePara CheckBalancePara
	ActionBase
}

type CheckBalancePara struct {
	Address             common.Address
	Validators          []common.Address
	NetStake            *big.Int
	CheckCommission     bool
	CommissionValidator common.Address
}

func ReadOnly(methodName string) bool {
	return methodName[0:3] == "get"
}

func getMethodType(methodName string) MethodType {
	switch {
	case methodName == "checkBalance":
		return CHECK_BALANCE
	case methodName[0:3] == "get":
		return QUERY
	default:
		return TX
	}
}

func (c *RawCase) Pack() (Case, error) {
	var Nonce_Map = make(map[common.Address]uint64)
	var res = Case{
		index: c.Index,
	}
	for i := 0; i < len(c.Actions); i++ {
		rawAction := c.Actions[i]
		var action Action
		var err error
		switch getMethodType(rawAction.MethodName) {
		case QUERY, CHECK_BALANCE:
			action, err = rawAction.Pack(0)
		default:
			sender := rawAction.Sender.ToAddress()
			nonce, ok := Nonce_Map[sender]
			if !ok {
				Nonce_Map[sender] = 1
				nonce = 0
			} else {
				Nonce_Map[sender] += 1
			}
			action, err = rawAction.Pack(nonce)
		}
		if err != nil {
			return Case{}, err
		}
		res.actions = append(res.actions, action)
	}
	return res, nil
}

func (a *RawAction) Pack(nonce uint64) (Action, error) {
	switch getMethodType(a.MethodName) {
	case CHECK_BALANCE:
		return &CheckBalance{
			ActionBase:          a.ActionBase,
			Address:             a.CheckBalancePara.Address,
			Validators:          a.CheckBalancePara.Validators,
			NetStake:            a.CheckBalancePara.NetStake,
			CheckCommission:     a.CheckBalancePara.CheckCommission,
			CommissionValidator: a.CheckBalancePara.CommissionValidator,
		}, nil
	case QUERY:
		data, err := a.Input.Encode()
		if err != nil {
			return nil, err
		}
		request := ethereum.CallMsg{To: &NODE_MANAGER_CONTRACT, Data: data}
		return &Query{
			ActionBase: a.ActionBase,
			Request:    request,
			Assertions: a.Assertions,
		}, nil
	default:
		data, err := a.Input.Encode()
		if err != nil {
			return nil, err
		}
		signKey := a.Sender.PrivateKey()
		log.Info("Packing tx", "sender", a.Sender.ToAddress().Hex(), "index_1", a.Sender.Index_1, "index_2", a.Sender.Index_2)
		tx := types.NewTransaction(nonce, NODE_MANAGER_CONTRACT, common.Big0, DEFAULT_GAS_LIMIT, DEFAULT_GAS_PRICE, data)
		signer := types.LatestSignerForChainID(ZION_CHAINID)
		tx, err = types.SignTx(tx, signer, signKey)
		if err != nil {
			return nil, err
		}
		return &SendTx{
			ActionBase:    a.ActionBase,
			Tx:            tx,
			ShouldSucceed: a.ShouldSucceed,
		}, nil

	}
}
