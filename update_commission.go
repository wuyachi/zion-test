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

func (u *UpdateCommissionParser) ParseInput(input string) error {
	param := &node_manager.UpdateCommissionParam{}
	u.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format input[%s]", input)
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		return err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()

	commission, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid commission: %s, err: %v", parts[1], err)
	}
	param.Commission = big.NewInt(commission)

	return nil
}

func (u *UpdateCommissionParser) ParseAssertion(input string) error {
	return nil
}
