package excel

import (
	"main/common"
	"strings"
)

type CreateValidatorParam struct {
	ConsensusAddress common.HDAddress
	SignerAddress    common.HDAddress
	ProposalAddress  common.HDAddress
	Commission       uint64
	InitStake        uint64
	Desc             string
}

type CreateValidatorComposer struct {
	rawAction *common.RawAction
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
