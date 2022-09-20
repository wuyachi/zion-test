package main

import (
	"fmt"
	"main/proposal_manager"
	"math/big"
	"strings"
)

type VoteProposalParser struct {
	rawAction *RawAction
}

func (c *VoteProposalParser) ParseInput(input string) error {
	param := &proposal_manager.VoteProposalParam{}
	c.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 1 {
		return fmt.Errorf("invalid format input[%s]", input)
	}

	id, ok := new(big.Float).SetString(parts[0])
	if !ok {
		return fmt.Errorf("invalid id: %s", parts[0])
	}
	idUint, _ := id.Uint64()
	param.ID = new(big.Int).SetUint64(idUint)

	return nil
}

func (c *VoteProposalParser) ParseAssertion(input string) error {
	return nil
}
