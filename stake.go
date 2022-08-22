package main

import (
	"fmt"
	"main/node_manager"
	"math/big"
	"strconv"
	"strings"
)

type StakeParser struct {
	rawAction *RawAction
}

func (s *StakeParser) ParseInput(input string) (Param, error) {
	param := &node_manager.StakeParam{}

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

	amount, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid stake amount: %s", parts[3])
		return nil, err
	}
	param.Amount = big.NewInt(amount)

	return param, nil
}

func (s *StakeParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
