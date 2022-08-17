package main

import (
	"strings"
)

type CreateValidatorParam struct {
	ConsensusAddress HDAddress
	SignerAddress    HDAddress
	ProposalAddress  HDAddress
	Commission       uint64
	InitStake        uint64
	Desc             string
}

type CreateValidatorComposer struct {
	rawAction *RawAction
}

func (c CreateValidatorComposer) compose() error {
	return nil
}

func parseCreateValidatorParam(input string) (param CreateValidatorParam, err error) {
	param = CreateValidatorParam{}
	paramStrs := strings.Split(input, ";")
	_ = paramStrs[0]
	// todo

	return
}

func composeCreateValidatorInput() (input []byte, err error) {
	//todo
	return
}
