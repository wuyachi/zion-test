package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/polynetwork/bridge-common/log"
	"main/proposal_manager"

	"main/base"
	"main/node_manager"
	"reflect"
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
			fmt.Printf("Assert decodeResult err:%s\n", err)
			return err
		}
		err = assertField(s, assertion.AssertType, assertion.FieldValues)
		if err != nil {
			fmt.Printf("Assert assertField err:%s\n", err)
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
				return fmt.Errorf("%s.%s should equal, expect %s, but got %s", result.Type().Name(), field, expect.Interface(), val.Interface())
			}
		case Assert_Element_Not_Equal:
			if checkType(val, expect) != 1 {
				return fmt.Errorf("%s.%s.Assert_Element_Not_Equal receive invalid type value, expect %s, but got %s", result.Type().Name(), field, expect.Type().Name(), val.Type().Name())
			}
			if equal(val, expect) {
				return fmt.Errorf("%s.%s should not equal, expect %s, but got %s", result.Type().Name(), field, expect.Interface(), val.Interface())
			}
		case Assert_Element_Contain:
			if checkType(val, expect) != 2 {
				return fmt.Errorf("%s.%s.Assert_Element_Contain receive invalid type value, expect element of %s, but got %s", result.Type().Name(), field, expect.Type().Name(), val.Type().Name())
			}
			if !contain(val, expect) {
				return fmt.Errorf("%s.%s should contain %s, but not", result.Type().Name(), field, expect.Interface())
			}
		case Assert_Element_Not_Contain:
			if checkType(val, expect) != 2 {
				return fmt.Errorf("%s.%s.Assert_Element_Not_Contain receive invalid type value, expect element of %s, but got %s", result.Type().Name(), field, expect.Type().Name(), val.Type().Name())
			}
			if contain(val, expect) {
				return fmt.Errorf("%s.%s should not contain %s, but still contain", result.Type().Name(), field, expect.Interface())
			}
		default:
			return fmt.Errorf("unknown AssertType: %d ", AssertType)
		}
	}

	return nil
}

func decodeResult(result []byte, methodName string) (res reflect.Value, err error) {
	var unpacked []interface{}
	zionContract := getZionContractAddress(methodName)
	switch zionContract {
	case NODE_MANAGER_CONTRACT:
		unpacked, err = node_manager.ABI.Unpack(methodName, result)
		if err != nil {
			return reflect.Value{}, err
		}
	case PROPOSAL_MANAGER_CONTRACT:
		unpacked, err = proposal_manager.ABI.Unpack(methodName, result)
		if err != nil {
			return reflect.Value{}, err
		}
	default:
		err = fmt.Errorf("undefined method:%s", methodName)
		return reflect.Value{}, err
	}

	result = *abi.ConvertType(unpacked[0], new([]byte)).(*[]byte)
	m, err := getMethodResult(methodName)
	if err != nil {
		return reflect.Value{}, err
	}
	err = rlp.DecodeBytes(result, m)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("fail to decode return value: %v %x", err, result)
	}
	log.Info("query", "method", methodName, "result", fmt.Sprintf("%+v", m))
	return reflect.ValueOf(m).Elem(), nil
}

func getMethodResult(methodName string) (interface{}, error) {
	switch methodName {
	case base.MethodGetAccumulatedCommission:
		return &node_manager.AccumulatedCommission{}, nil
	case base.MethodGetAllValidators:
		return &node_manager.AllValidators{}, nil
	case base.MethodGetCommunityInfo:
		return &node_manager.CommunityInfo{}, nil
	case base.MethodGetCurrentEpochInfo:
		return &node_manager.EpochInfo{}, nil
	case base.MethodGetEpochInfo:
		return &node_manager.EpochInfo{}, nil
	case base.MethodGetGlobalConfig:
		return &node_manager.GlobalConfig{}, nil
	case base.MethodGetOutstandingRewards:
		return &node_manager.OutstandingRewards{}, nil
	case base.MethodGetStakeInfo:
		return &node_manager.StakeInfo{}, nil
	case base.MethodGetStakeStartingInfo:
		return &node_manager.StakeStartingInfo{}, nil
	case base.MethodGetTotalPool:
		return &node_manager.TotalPool{}, nil
	case base.MethodGetUnlockingInfo:
		return &node_manager.UnlockingInfo{}, nil
	case base.MethodGetValidator:
		return &node_manager.Validator{}, nil
	case base.MethodGetValidatorAccumulatedRewards:
		return &node_manager.ValidatorAccumulatedRewards{}, nil
	case base.MethodGetValidatorOutstandingRewards:
		return &node_manager.ValidatorOutstandingRewards{}, nil
	case base.MethodGetValidatorSnapshotRewards:
		return &node_manager.ValidatorSnapshotRewards{}, nil
	case base.MethodGetStakeRewards:
		return &node_manager.ValidatorOutstandingRewards{}, nil
	case base.MethodGetProposal:
		return &proposal_manager.Proposal{}, nil
	case base.MethodGetProposalList:
		return &proposal_manager.ProposalList{}, nil
	case base.MethodGetConfigProposalList:
		return &proposal_manager.ConfigProposalList{}, nil
	case base.MethodGetCommunityProposalList:
		return &proposal_manager.CommunityProposalList{}, nil

	default:
		err := fmt.Errorf("getMethodResult undefined method: %s", methodName)
		return nil, err
	}
}
