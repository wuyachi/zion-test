package main

import (
	"fmt"
	"main/base"
	"main/node_manager"
	"math/big"
	"strconv"
	"strings"
)

type StakeParser struct {
	rawAction *RawAction
}

func (s *StakeParser) ParseInput(input string) error {
	param := &node_manager.StakeParam{}
	s.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format input[%s]", input)
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		return err
	}

	param.ConsensusAddress = consensusHdAddress.ToAddress()

	amount, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid stake amount: %s, err: %v", parts[1], err)
	}
	s.rawAction.Amount = new(big.Int).Mul(big.NewInt(amount), base.ZionPrecision)

	return nil
}

func (s *StakeParser) ParseAssertion(input string) error {
	return nil
}
