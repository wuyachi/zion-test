package main

import (
	"main/node_manager"
)

type ChangeEpochParser struct {
	rawAction *RawAction
}

func (c *ChangeEpochParser) ParseInput(input string) (Param, error) {
	param := &node_manager.ChangeEpochParam{}
	return param, nil
}

func (c *ChangeEpochParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
