package main

import (
	"main/node_manager"
)

type WithdrawParser struct {
	rawAction *RawAction
}

func (w *WithdrawParser) ParseInput(input string) (Param, error) {
	param := &node_manager.WithdrawParam{}
	return param, nil
}

func (w *WithdrawParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
