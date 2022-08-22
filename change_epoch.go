package main

import (
	"main/node_manager"
)

type ChangeEpochParser struct {
	rawAction *RawAction
}

func (c *ChangeEpochParser) ParseInput(input string) error {
	c.rawAction.Input = &node_manager.ChangeEpochParam{}
	return nil
}

func (c *ChangeEpochParser) ParseAssertion(input string) error {
	return nil
}
