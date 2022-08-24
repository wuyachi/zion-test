package main

import (
	"fmt"
	"main/base"
	"main/node_manager"
	"strings"
)

type GetCurrentEpochInfoParser struct {
	rawAction *RawAction
}

func (g *GetCurrentEpochInfoParser) ParseInput(input string) error {
	g.rawAction.Input = &node_manager.GetCurrentEpochInfoParam{}
	return nil
}

func (g *GetCurrentEpochInfoParser) ParseAssertion(input string) error {
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
	case "Validators":
		assertion := Assertion{}
		assertion.AssertType = assertType
		assertion.MethodName = base.MethodGetCurrentEpochInfo
		fieldValues := make([]FieldValue, 0)
		for _, value := range values {
			fieldValue := FieldValue{}
			fieldValue.Field = field
			hdAddress, e := parseAddress(value)
			if e != nil {
				return e
			}
			address := hdAddress.ToAddress()
			fieldValue.Value = address
			fieldValues = append(fieldValues, fieldValue)
		}
		assertion.FieldValues = fieldValues
		g.rawAction.Assertions = append(g.rawAction.Assertions, assertion)
	default:
		return fmt.Errorf("undefined assertion field: %s", field)

	}

	return nil
}
