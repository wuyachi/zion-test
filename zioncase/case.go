package zioncase

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
)

type Case struct {
	actions []Action
}

type Action interface {
	Run() error
}

type ActionBase struct {
	Block        uint64
	ShouldBefore uint64
	EpochId      uint64
}

type SendTx struct {
	ActionBase
	Tx            []byte
	ShouldSucceed bool
}

type Query struct {
	ActionBase
	Request        []byte
	ExpectedResult []byte
}

type Param interface {
	Encode() ([]byte, error)
}

const TEST_MM = "test test test test test test test test test test test test"

type HDAddress struct {
	Index_1 uint64
	Index_2 uint64
}

func (hd *HDAddress) ToAddress() common.Address    { return common.Address{} }
func (hf *HDAddress) PrivateKey() ecdsa.PrivateKey { return ecdsa.PrivateKey{} }

type RawCase struct {
	actions []RawAction
}

type RawAction struct {
	MethodName    string
	RawInput      string
	Input         Param
	ShouldSucceed bool
	Result        []byte
	Sender        HDAddress
	Options       ActionBase
}
