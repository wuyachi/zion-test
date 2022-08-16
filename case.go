package zion_test

import (
	"github.com/ethereum/go-ethereum/common"
)


type Case struct {
	actions []Action
}

type Action interface {
	Run() error
}

type ActionBase struct {
	Block 			uint64
	ShouldBefore    uint64
	Sender          common.Address
}

type SendTx struct {
	ActionBase
	Tx []byte
	ShouldSucceed bool
}

func (a *SendTx) Run() error {
	return nil
}

type Query struct {
	ActionBase
	Request []byte
	ExpectedResult []byte
}

func (a *Query) Run() error {
	return nil
}