package main

import (
	"fmt"
	"main/base"
	"main/node_manager"
	"math/big"
	"strconv"
	"strings"
)

type CreateValidatorParser struct {
	rawAction *RawAction
}

func (c *CreateValidatorParser) ParseInput(input string) error {
	param := &node_manager.CreateValidatorParam{}
	c.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 6 {
		return fmt.Errorf("invalid format input[%s]", input)
	}
	consensusHdAddress, err := parseAddress(parts[0])
	if err != nil {
		return err
	}
	signerHdAddress, err := parseAddress(parts[1])
	if err != nil {
		return err
	}
	proposalHdAddress, err := parseAddress(parts[2])
	if err != nil {
		return err
	}
	param.ConsensusAddress = consensusHdAddress.ToAddress()
	param.SignerAddress = signerHdAddress.ToAddress()
	param.ProposalAddress = proposalHdAddress.ToAddress()

	commission, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid commission: %s, err: %v", parts[3], err)
	}
	initStake, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid initStake: %s, err: %v", parts[4], err)
	}
	param.Commission = big.NewInt(commission)

	c.rawAction.Amount = new(big.Int).Mul(big.NewInt(initStake), base.ZionPrecision)
	param.Desc = parts[5]

	return nil
}

func (c *CreateValidatorParser) ParseAssertion(input string) error {
	return nil
}
