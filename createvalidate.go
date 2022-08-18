package main

type CreateValidatorParam struct {
	ConsensusAddress HDAddress
	SignerAddress    HDAddress
	ProposalAddress  HDAddress
	Commission       uint64
	InitStake        uint64
	Desc             string
}

type CreateValidatorParser struct {
	rawAction *RawAction
}

func (c CreateValidatorParser) parseInput() error {
	return nil
	//todo
}

func (c CreateValidatorParser) parseAssertion() error {
	return nil
	//todo
}
