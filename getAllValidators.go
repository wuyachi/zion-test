package main

import (
	"fmt"
	"main/base"
	"main/node_manager"
	"strings"
)

type GetAllValidatorsParser struct {
	rawAction *RawAction
}

func (g *GetAllValidatorsParser) ParseInput(input string) (Param, error) {
	param := &node_manager.GetAllValidatorsParam{}
	return param, nil
}

func (g *GetAllValidatorsParser) ParseAssertion(input string) (assertions []Assertion, err error) {
	if input == "nil" {
		return nil, nil
	}

	parts := strings.Split(input, ";")

	field := parts[1]
	assertType, err := formatAssertType(parts[0])
	if err != nil {
		err = fmt.Errorf("invalid format, err: %v", err)
		return
	}

	assertion := Assertion{}
	assertion.AssertType = assertType
	assertion.MethodName = base.MethodGetAllValidators

	values := parts[2:]
	fieldValues := make([]FieldValue, 0)

	switch field {
	case "AllValidators":
		for _, value := range values {
			fieldValue := FieldValue{}
			fieldValue.Field = field
			hdAddress, e := parseAddress(value)
			if e != nil {
				err = fmt.Errorf("parse validator address failed, addressTag: %s, err: %v", value, e)
				return
			}
			address := hdAddress.ToAddress()
			fieldValue.Value = address
			fieldValues = append(fieldValues, fieldValue)
		}
	default:
		err = fmt.Errorf("undefined assertion field: %s", field)
		return
	}
	assertion.FieldValues = fieldValues

	return
}
