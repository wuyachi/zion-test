package main

import (
	"main/node_manager"
)

type WithdrawParser struct {
	rawAction *RawAction
}

func (w *WithdrawParser) ParseInput(input string) error {
	w.rawAction.Input = &node_manager.WithdrawParam{}
	return nil
}

func (w *WithdrawParser) ParseAssertion(input string) error {
	return nil
}
