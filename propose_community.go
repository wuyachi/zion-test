package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"main/node_manager"
	"main/proposal_manager"
	"math/big"
	"strings"
)

type ProposeCommunityParser struct {
	rawAction *RawAction
}

func (c *ProposeCommunityParser) ParseInput(input string) error {
	param := &proposal_manager.ProposeCommunityParam{}
	c.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format input[%s]", input)
	}

	communityRate, ok := new(big.Int).SetString(parts[0], 10)
	if !ok {
		return fmt.Errorf("invalid communityRate: %s", parts[0])
	}
	communityAddress, err := parseAddress(parts[1])
	if err != nil {
		return err
	}
	communityInfo := &node_manager.CommunityInfo{communityRate, communityAddress.ToAddress()}
	content, err := rlp.EncodeToBytes(communityInfo)
	if err != nil {
		return fmt.Errorf("rlp.EncodeToBytes(communityInfo) err: %v", err)
	}
	param.Content = content

	return nil
}

func (c *ProposeCommunityParser) ParseAssertion(input string) error {
	return nil
}
