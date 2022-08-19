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

func (g *GetCurrentEpochInfoParser) ParseInput(input string) (Param, error) {
	param := &node_manager.GetCurrentEpochInfoParam{}
	return param, nil
}

func (g *GetCurrentEpochInfoParser) ParseAssertion(input string) (assertions []Assertion, err error) {
	if input == "nil" {
		return nil, nil
	}

	parts := strings.Split(input, ";")

	field := parts[1]
	assertType, err := formatAssertType(parts[0])
	if err != nil {
		err = fmt.Errorf("ParseAssertion formatAssertType failed. err:%v", err)
		return
	}

	assertion := Assertion{}
	assertion.AssertType = assertType
	assertion.MethodName = base.MethodGetCurrentEpochInfo

	values := parts[2:]
	fieldValues := make([]FieldValue, 0)

	switch field {
	case "Validators":
		for _, value := range values {
			fieldValue := FieldValue{}
			fieldValue.Field = field
			hdAddress, e := parseAddress(value)
			if e != nil {
				err = fmt.Errorf("parse validator address failed. addressTag:%s, err:%v", value, e)
				return
			}
			address := hdAddress.ToAddress()
			fieldValue.Value = address
			fieldValues = append(fieldValues, fieldValue)
		}
	default:
		err = fmt.Errorf("ParseAssertion undefined field=%s", field)
		return
	}
	assertion.FieldValues = fieldValues

	return
}
