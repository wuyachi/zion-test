package main

import (
	"fmt"
	"main/node_manager"
	"math/big"
	"strconv"
	"strings"
)

type UpdateCommissionParser struct {
	rawAction *RawAction
}

func (u *UpdateCommissionParser) ParseInput(input string) (Param, error) {
	param := &node_manager.UpdateCommissionParam{}

	parts := strings.Split(input, ";")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format input[%s]", input)
		return nil, err
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		err = fmt.Errorf("parse consensusAddress failed, input: %s", input)
		return nil, err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()

	commission, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid commission: %s", parts[1])
		return nil, err
	}
	param.Commission = big.NewInt(commission)

	return param, nil
}

func (u *UpdateCommissionParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
