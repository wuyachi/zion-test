package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/accounts/abi"

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
	unpacked, err := node_manager.ABI.Unpack(methodName, result)
	if err != nil {
		return reflect.Value{}, err
	}
	result = *abi.ConvertType(unpacked[0], new([]byte)).(*[]byte)
	m, ok := MethodResultMap[methodName]
	if !ok {
		return reflect.Value{}, fmt.Errorf("unknown method name:%s", methodName)
	}
	err = rlp.DecodeBytes(result, &m)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("fail to decode return value: %v %x", err, result)
	}
	return reflect.ValueOf(m), nil
}

var MethodResultMap = map[string]interface{}{
	base.MethodGetAccumulatedCommission:       node_manager.AccumulatedCommission{},
	base.MethodGetAllValidators:               node_manager.AllValidators{},
	base.MethodGetCommunityInfo:               node_manager.CommunityInfo{},
	base.MethodGetCurrentEpochInfo:            node_manager.EpochInfo{},
	base.MethodGetEpochInfo:                   node_manager.EpochInfo{},
	base.MethodGetGlobalConfig:                node_manager.GlobalConfig{},
	base.MethodGetOutstandingRewards:          node_manager.OutstandingRewards{},
	base.MethodGetStakeInfo:                   node_manager.StakeInfo{},
	base.MethodGetStakeStartingInfo:           node_manager.StakeStartingInfo{},
	base.MethodGetTotalPool:                   node_manager.TotalPool{},
	base.MethodGetUnlockingInfo:               node_manager.UnlockingInfo{},
	base.MethodGetValidator:                   node_manager.Validator{},
	base.MethodGetValidatorAccumulatedRewards: node_manager.ValidatorAccumulatedRewards{},
	base.MethodGetValidatorOutstandingRewards: node_manager.ValidatorOutstandingRewards{},
	base.MethodGetValidatorSnapshotRewards:    node_manager.ValidatorSnapshotRewards{},
}
