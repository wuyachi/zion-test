package main

import (
	"main/node_manager"
)

type EndBlockParser struct {
	rawAction *RawAction
}

func (b *EndBlockParser) ParseInput(input string) (Param, error) {
	param := &node_manager.EndBlockParam{}
	return param, nil
}

func (b *EndBlockParser) ParseAssertion(input string) ([]Assertion, error) {
	return nil, nil
}
