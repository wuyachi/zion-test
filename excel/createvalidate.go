package excel

import (
	"strings"
	"zion-test/zioncase"
)

type CreateValidatorParam struct {
	ConsensusAddress zioncase.HDAddress
	SignerAddress    zioncase.HDAddress
	ProposalAddress  zioncase.HDAddress
	Commission       uint64
	InitStake        uint64
	Desc             string
}

type CreateValidatorComposer struct {
	rawAction *zioncase.RawAction
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
