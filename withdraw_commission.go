package main

import (
	"fmt"
	"main/node_manager"
	"strings"
)

type WithdrawCommissionParser struct {
	rawAction *RawAction
}

func (w *WithdrawCommissionParser) ParseInput(input string) error {
	param := &node_manager.WithdrawCommissionParam{}
	w.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 1 {
		return fmt.Errorf("invalid format input[%s]", input)
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		return err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()

	return nil
}

func (w *WithdrawCommissionParser) ParseAssertion(input string) error {
	return nil
}
