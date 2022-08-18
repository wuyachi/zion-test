package main

import (
	"fmt"
	"math/big"
	"reflect"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type AssertType uint8

const (
	Assert_Element_Contain     AssertType = 0
	Assert_Element_Not_Contain AssertType = 1
	Assert_Element_Equal       AssertType = 2
	Assert_Element_Not_Equal   AssertType = 3
)

type Assertion struct {
	AssertType
	MethodName  string
	FieldValues []FieldValue
}

type FieldValue struct {
	Field string
	Value interface{}
}

func Assert(result []byte, assertions []Assertion) error {
	for i := 0; i < len(assertions); i++ {
		assertion := assertions[i]
		s, err := decodeResult(result, assertion.MethodName)
		if err != nil {
			return err
		}
		err = assertField(s, assertion.AssertType, assertion.FieldValues)
		if err != nil {
			return err
		}
	}
	return nil
}

// return 0 if x.Type != y.Type && x.Type != []y.Type
// return 1 if x.Type == y.Type
// return 2 if x.Type == []y.Type
func checkType(x, y reflect.Value) int8 {
	if x.Type() == y.Type() {
		return 1
	}
	if x.Type() == reflect.SliceOf(y.Type()) {
		return 2
	}
	return 0
}

// will panic if x is not slice/array
// return if x contains y
func contain(x, y reflect.Value) bool {
	for i := 0; i < x.Len(); i++ {
		if equal(x.Index(i), y) {
			return true
		}
	}
	return false
}

func equal(x, y reflect.Value) bool {
	return reflect.DeepEqual(x.Interface(), y.Interface())
}

func assertField(result reflect.Value, AssertType AssertType, fieldValues []FieldValue) error {
	for i := 0; i < len(fieldValues); i++ {
		field := fieldValues[i].Field
		expect := reflect.ValueOf(fieldValues[i].Value)
		val := result.FieldByName(field)
		switch AssertType {
		case Assert_Element_Equal:
			if checkType(val, expect) != 1 {
				return fmt.Errorf("%s.%s.Assert_Element_Equal receive invalid type value, expect %s, but got %s", result.Type().Name(), field, expect.Type().Name(), val.Type().Name())
			}
			if !equal(val, expect) {
				return fmt.Errorf("%s.%s should equal, expect %s, but got %s", result.Type().Name(), field, expect.String(), val.String())
			}
		case Assert_Element_Not_Equal:
			if checkType(val, expect) != 1 {
				return fmt.Errorf("%s.%s.Assert_Element_Not_Equal receive invalid type value, expect %s, but got %s", result.Type().Name(), field, expect.Type().Name(), val.Type().Name())
			}
			if equal(val, expect) {
				return fmt.Errorf("%s.%s should not equal, expect %s, but got %s", result.Type().Name(), field, expect.String(), val.String())
			}
		case Assert_Element_Contain:
			if checkType(val, expect) != 2 {
				return fmt.Errorf("%s.%s.Assert_Element_Contain receive invalid type value, expect element of %s, but got %s", result.Type().Name(), field, expect.Type().Name(), val.Type().Name())
			}
			if !contain(val, expect) {
				return fmt.Errorf("%s.%s should contain %s, but not", result.Type().Name(), field, expect.String())
			}
		case Assert_Element_Not_Contain:
			if checkType(val, expect) != 2 {
				return fmt.Errorf("%s.%s.Assert_Element_Not_Contain receive invalid type value, expect element of %s, but got %s", result.Type().Name(), field, expect.Type().Name(), val.Type().Name())
			}
			if contain(val, expect) {
				return fmt.Errorf("%s.%s should not contain %s, but still contain", result.Type().Name(), field, expect.String())
			}
		default:
			return fmt.Errorf("unknown AssertType: %d ", AssertType)
		}
	}

	return nil
}

func decodeResult(result []byte, methodName string) (reflect.Value, error) {
	m, ok := MethodResultMap[methodName]
	if !ok {
		return reflect.Value{}, fmt.Errorf("unknown method name:%s", methodName)
	}
	err := rlp.DecodeBytes(result, m)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("fail to decode return value")
	}
	return reflect.ValueOf(m), nil
}

var MethodResultMap = map[string]interface{}{
	"getAccumulatedCommission":       AccumulatedCommission{},
	"getAllValidators":               AllValidators{},
	"getCommunityInfo":               CommunityInfo{},
	"getCurrentEpochInfo":            EpochInfo{},
	"getEpochInfo":                   EpochInfo{},
	"getGlobalConfig":                GlobalConfig{},
	"getOutstandingRewards":          OutstandingRewards{},
	"getStakeInfo":                   StakeInfo{},
	"getStakeStartingInfo":           StakeStartingInfo{},
	"getTotalPool":                   TotalPool{},
	"getUnlockingInfo":               UnlockingInfo{},
	"getValidator":                   Validator{},
	"getValidatorAccumulatedRewards": ValidatorAccumulatedRewards{},
	"getValidatorOutstandingRewards": ValidatorOutstandingRewards{},
	"getValidatorSnapshotRewards":    ValidatorSnapshotRewards{},
}

type Dec struct {
	I *big.Int
}

type LockStatus uint8

const (
	Unspecified LockStatus = 0
	Unlock      LockStatus = 1
	Lock        LockStatus = 2
	Remove      LockStatus = 3
)

type AllValidators struct {
	AllValidators []common.Address
}

type Validator struct {
	StakeAddress     common.Address
	ConsensusAddress common.Address
	SignerAddress    common.Address
	ProposalAddress  common.Address
	Commission       *Commission
	Status           LockStatus
	Jailed           bool
	UnlockHeight     *big.Int
	TotalStake       Dec
	SelfStake        Dec
	Desc             string
}

type Commission struct {
	Rate         Dec
	UpdateHeight *big.Int
}

type GlobalConfig struct {
	MaxCommissionChange   *big.Int
	MinInitialStake       *big.Int
	MinProposalStake      *big.Int
	BlockPerEpoch         *big.Int
	ConsensusValidatorNum uint64
	VoterValidatorNum     uint64
}

type StakeInfo struct {
	StakeAddress  common.Address
	ConsensusAddr common.Address
	Amount        Dec
}

type UnlockingInfo struct {
	StakeAddress   common.Address
	UnlockingStake []*UnlockingStake
}

type UnlockingStake struct {
	Height           *big.Int
	CompleteHeight   *big.Int
	ConsensusAddress common.Address
	Amount           Dec
}

type EpochInfo struct {
	ID          *big.Int
	Validators  []common.Address
	Signers     []common.Address
	Voters      []common.Address
	Proposers   []common.Address
	StartHeight *big.Int
	EndHeight   *big.Int
}

type AccumulatedCommission struct {
	Amount Dec
}

type ValidatorAccumulatedRewards struct {
	Rewards Dec
	Period  uint64
}

type ValidatorOutstandingRewards struct {
	Rewards Dec
}

type OutstandingRewards struct {
	Rewards Dec
}

type ValidatorSnapshotRewards struct {
	AccumulatedRewardsRatio Dec
	ReferenceCount          uint64
}

type StakeStartingInfo struct {
	StartPeriod uint64
	Stake       Dec
	Height      *big.Int
}

type AddressList struct {
	List []common.Address
}

type ConsensusSign struct {
	Method string
	Input  []byte
	hash   atomic.Value
}

type CommunityInfo struct {
	CommunityRate    *big.Int
	CommunityAddress common.Address
}

type TotalPool struct {
	TotalPool Dec
}
