package main

import (
	"fmt"
	"main/node_manager"
	"strings"
)

type CancelValidatorParser struct {
	rawAction *RawAction
}

func (c *CancelValidatorParser) ParseInput(input string) error {
	param := &node_manager.CancelValidatorParam{}
	c.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 1 {
		err := fmt.Errorf("invalid format input[%s]", input)
		return err
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		err = fmt.Errorf("parse consensusAddress failed, input: %s", input)
		return err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()

	return nil
}

func (c *CancelValidatorParser) ParseAssertion(input string) error {
	return nil
}
