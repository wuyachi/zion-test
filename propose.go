package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"main/base"
	"main/node_manager"
	"main/proposal_manager"
	"math/big"
	"strconv"
	"strings"
)

type ProposeParser struct {
	rawAction *RawAction
}

func (c *ProposeParser) ParseInput(input string) error {
	param := &proposal_manager.ProposeParam{}
	c.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 6 {
		return fmt.Errorf("invalid format input[%s]", input)
	}

	maxCommissionChange, ok := new(big.Int).SetString(parts[0], 10)
	if !ok {
		return fmt.Errorf("invalid maxCommissionChange: %s", parts[0])
	}

	minInitialStake, ok := new(big.Int).SetString(parts[1], 10)
	if !ok {
		return fmt.Errorf("invalid minInitialStake: %s", parts[1])
	}
	minInitialStake = new(big.Int).Mul(minInitialStake, base.ZionPrecision)

	minProposalStake, ok := new(big.Int).SetString(parts[2], 10)
	if !ok {
		return fmt.Errorf("invalid minProposalStake: %s", parts[2])
	}
	minProposalStake = new(big.Int).Mul(minProposalStake, base.ZionPrecision)

	blockPerEpoch, ok := new(big.Int).SetString(parts[3], 10)
	if !ok {
		return fmt.Errorf("invalid minProposalStake: %s", parts[3])
	}

	consensusValidatorNum, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid consensusValidatorNum: %s, err: %v", parts[4], err)
	}

	voterValidatorNum, err := strconv.ParseUint(parts[5], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid voterValidatorNum: %s, err: %v", parts[5], err)
	}

	globalConfig := &node_manager.GlobalConfig{
		maxCommissionChange,
		minInitialStake,
		minProposalStake,
		blockPerEpoch,
		consensusValidatorNum,
		voterValidatorNum,
	}
	content, err := rlp.EncodeToBytes(globalConfig)
	if err != nil {
		return fmt.Errorf("rlp.EncodeToBytes(globalConfig) err: %v", err)
	}
	param.Content = content

	return nil
}

func (c *ProposeParser) ParseAssertion(input string) error {
	return nil
}
