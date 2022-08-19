package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type CreateValidatorParser struct {
	rawAction *RawAction
}

func (c *CreateValidatorParser) ParseInput(input string) (Param, error) {
	param := &CreateValidatorParam{}

	parts := strings.Split(input, ";")
	if len(parts) != 6 {
		err := fmt.Errorf("invalid format input[%s]", input)
		return nil, err
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		err = fmt.Errorf("parse consensusAddress failed. input=%s", input)
		return nil, err
	}
	signerHdAddress, err := parseAddress(parts[1])
	if err != nil {
		err = fmt.Errorf("parse signerAddress failed. input=%s", input)
		return nil, err
	}
	proposalHdAddress, err := parseAddress(parts[2])
	if err != nil {
		err = fmt.Errorf("parse proposalAddress failed. input=%s", input)
		return nil, err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()
	param.SignerAddress = signerHdAddress.ToAddress()
	param.ProposalAddress = proposalHdAddress.ToAddress()

	commission, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid commission:%s", parts[3])
		return nil, err
	}
	initStake, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid initStake:%s", parts[4])
		return nil, err
	}
	param.Commission = big.NewInt(commission)
	param.InitStake = big.NewInt(initStake)
	param.Desc = parts[5]

	return param, nil
}

func (c *CreateValidatorParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
