package main

import (
	"fmt"
	"main/node_manager"
	"strings"
)

type UpdateValidatorParser struct {
	rawAction *RawAction
}

func (u *UpdateValidatorParser) ParseInput(input string) error {
	param := &node_manager.UpdateValidatorParam{}
	u.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 4 {
		return fmt.Errorf("invalid format input[%s]", input)
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		return err
	}
	signerAddress, err := parseAddress(parts[1])
	if err != nil {
		return err
	}
	proposalHdAddress, err := parseAddress(parts[2])
	if err != nil {
		return err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()
	param.SignerAddress = signerAddress.ToAddress()
	param.ProposalAddress = proposalHdAddress.ToAddress()
	param.Desc = parts[3]

	return nil
}

func (u *UpdateValidatorParser) ParseAssertion(input string) error {
	return nil
}
