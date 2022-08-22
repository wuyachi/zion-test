package main

import (
	"main/node_manager"
)

type EndBlockParser struct {
	rawAction *RawAction
}

func (b *EndBlockParser) ParseInput(input string) error {
	b.rawAction.Input = &node_manager.EndBlockParam{}
	return nil
}

func (b *EndBlockParser) ParseAssertion(input string) error {
	return nil
}
