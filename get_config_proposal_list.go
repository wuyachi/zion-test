package main

import (
	"fmt"
	"main/base"
	"main/proposal_manager"
	"math/big"
	"strings"
)

type GetConfigProposalListParser struct {
	rawAction *RawAction
}

func (g *GetConfigProposalListParser) ParseInput(input string) error {
	g.rawAction.Input = &proposal_manager.GetConfigProposalListParam{}
	return nil
}

func (g *GetConfigProposalListParser) ParseAssertion(input string) error {
	if input == "nil" {
		return nil
	}

	parts := strings.Split(input, ";")

	field := parts[1]
	assertType, err := formatAssertType(parts[0])
	if err != nil {
		return err
	}

	values := parts[2:]
	switch field {
	case "ConfigProposalList":
		assertion := Assertion{}
		assertion.AssertType = assertType
		assertion.MethodName = base.MethodGetConfigProposalList
		fieldValues := make([]FieldValue, 0)
		for _, value := range values {
			fieldValue := FieldValue{}
			fieldValue.Field = field
			id, ok := new(big.Int).SetString(value, 10)
			if !ok {
				return fmt.Errorf("invalid id: %s", value)
			}
			fieldValue.Value = id
			fieldValues = append(fieldValues, fieldValue)
		}
		assertion.FieldValues = fieldValues
		g.rawAction.Assertions = append(g.rawAction.Assertions, assertion)
	default:
		return fmt.Errorf("undefined assertion field: %s", field)
	}

	return nil
}
