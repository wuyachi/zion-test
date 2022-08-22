package main

import (
	"fmt"
	"main/node_manager"
	"strings"
)

type UpdateValidatorParser struct {
	rawAction *RawAction
}

func (u *UpdateValidatorParser) ParseInput(input string) (Param, error) {
	param := &node_manager.UpdateValidatorParam{}

	parts := strings.Split(input, ";")
	if len(parts) != 3 {
		err := fmt.Errorf("invalid format input[%s]", input)
		return nil, err
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		err = fmt.Errorf("parse consensusAddress failed, input: %s", input)
		return nil, err
	}
	proposalHdAddress, err := parseAddress(parts[1])
	if err != nil {
		err = fmt.Errorf("parse proposalAddress failed, input: %s", input)
		return nil, err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()
	param.ProposalAddress = proposalHdAddress.ToAddress()

	param.Desc = parts[2]

	return param, nil
}

func (u *UpdateValidatorParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
