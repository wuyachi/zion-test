package main

import (
	"fmt"
	"main/node_manager"
	"strings"
)

type WithdrawStakeRewardsParser struct {
	rawAction *RawAction
}

func (w *WithdrawStakeRewardsParser) ParseInput(input string) (Param, error) {
	param := &node_manager.WithdrawStakeRewardsParam{}

	parts := strings.Split(input, ";")
	if len(parts) != 1 {
		err := fmt.Errorf("invalid format input[%s]", input)
		return nil, err
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		err = fmt.Errorf("parse consensusAddress failed, input: %s", input)
		return nil, err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()

	return param, nil
}

func (w *WithdrawStakeRewardsParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
