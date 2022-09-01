package main

import (
	"fmt"
	"main/base"
	"main/node_manager"
	"math/big"
	"strconv"
	"strings"
)

type GetStakeStartingInfoParser struct {
	rawAction *RawAction
}

func (g *GetStakeStartingInfoParser) ParseInput(input string) error {
	param := &node_manager.GetStakeStartingInfoParam{}
	g.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format input[%s]", input)
		return err
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		err = fmt.Errorf("parse consensusAddress failed, input: %s", input)
		return err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()

	stakeHdAddress, err := parseAddress(parts[1])
	if err != nil {
		err = fmt.Errorf("parse stakeAddress failed, input: %s", input)
		return err
	}
	param.StakeAddress = stakeHdAddress.ToAddress()

	return nil
}

func (g *GetStakeStartingInfoParser) ParseAssertion(input string) error {
	if input == "nil" {
		return nil
	}

	parts := strings.Split(input, ";")

	field := parts[1]
	assertType, err := formatAssertType(parts[0])
	if err != nil {
		return err
	}

	assertion := Assertion{}
	assertion.AssertType = assertType
	assertion.MethodName = base.MethodGetStakeStartingInfo
	fieldValues := make([]FieldValue, 0)
	switch field {
	case "Stake":
		amount, ok := new(big.Int).SetString(parts[2], 10)
		if !ok {
			return fmt.Errorf("invalid assertion format: %s", parts[2])
		}
		amount = amount.Mul(amount, base.ZionPrecision)
		amountDec := node_manager.NewDecFromBigInt(amount)

		fieldValue := FieldValue{}
		fieldValue.Field = field
		fieldValue.Value = amountDec
		fieldValues = append(fieldValues, fieldValue)
		assertion.FieldValues = fieldValues
		g.rawAction.Assertions = append(g.rawAction.Assertions, assertion)
	case "StartPeriod":
		startPeriod, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid assertion format: %s, err: %v", parts[2], err)
		}
		fieldValue := FieldValue{}
		fieldValue.Field = field
		fieldValue.Value = startPeriod
		fieldValues = append(fieldValues, fieldValue)
		assertion.FieldValues = fieldValues
		g.rawAction.Assertions = append(g.rawAction.Assertions, assertion)
	default:
		return fmt.Errorf("undefined assertion field: %s", field)
	}

	return nil
}
