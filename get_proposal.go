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

type GetProposalParser struct {
	rawAction *RawAction
}

func (g *GetProposalParser) ParseInput(input string) error {
	param := &proposal_manager.GetProposalParam{}
	g.rawAction.Input = param

	parts := strings.Split(input, ";")
	if len(parts) != 1 {
		return fmt.Errorf("invalid format input[%s]", input)
	}

	id, ok := new(big.Int).SetString(parts[0], 10)
	if !ok {
		return fmt.Errorf("invalid id: %s", parts[0])
	}
	param.ID = id

	return nil
}

func (g *GetProposalParser) ParseAssertion(input string) error {
	if input == "nil" {
		return nil
	}

	parts := strings.Split(input, ";")

	field := parts[1]
	assertType, err := formatAssertType(parts[0])
	if err != nil {
		return err
	}

	switch field {
	case "Content":
		assertion := Assertion{}
		assertion.AssertType = assertType
		assertion.MethodName = base.MethodGetProposal
		fieldValues := make([]FieldValue, 0)

		content, err := parseContent(parts[2:])
		if err != nil {
			return fmt.Errorf("invalid assertion format: %s", parts[2:])
		}

		fieldValue := FieldValue{}
		fieldValue.Field = field
		fieldValue.Value = content
		fieldValues = append(fieldValues, fieldValue)
		assertion.FieldValues = fieldValues
		g.rawAction.Assertions = append(g.rawAction.Assertions, assertion)
	default:
		return fmt.Errorf("undefined assertion field: %s", field)
	}

	return nil
}

func parseContent(parts []string) ([]byte, error) {
	maxCommissionChange, ok := new(big.Int).SetString(parts[0], 10)
	if !ok {
		return nil, fmt.Errorf("invalid maxCommissionChange: %s", parts[0])
	}

	minInitialStake, ok := new(big.Int).SetString(parts[1], 10)
	if !ok {
		return nil, fmt.Errorf("invalid minInitialStake: %s", parts[1])
	}
	minInitialStake = new(big.Int).Mul(minInitialStake, base.ZionPrecision)

	minProposalStake, ok := new(big.Int).SetString(parts[2], 10)
	if !ok {
		return nil, fmt.Errorf("invalid minProposalStake: %s", parts[2])
	}
	minProposalStake = new(big.Int).Mul(minProposalStake, base.ZionPrecision)

	blockPerEpoch, ok := new(big.Int).SetString(parts[3], 10)
	if !ok {
		return nil, fmt.Errorf("invalid minProposalStake: %s", parts[3])
	}

	consensusValidatorNum, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid consensusValidatorNum: %s, err: %v", parts[4], err)
	}

	voterValidatorNum, err := strconv.ParseUint(parts[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid voterValidatorNum: %s, err: %v", parts[5], err)
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
		return nil, fmt.Errorf("rlp.EncodeToBytes(globalConfig) err: %v", err)
	}
	return content, nil
}
