package main

import (
	"fmt"
	"math/big"
	"strings"
)

type CheckBalanceParser struct {
	rawAction *RawAction
}

func (c *CheckBalanceParser) ParseInput(input string) error {
	parts := strings.Split(input, ";")
	if len(parts) != 3 {
		return fmt.Errorf("invalid format input[%s]", input)
	}
	checkBalancePara := CheckBalancePara{}

	userHdAddress, err := parseAddress(parts[0])
	if err != nil {
		return err
	}
	checkBalancePara.Address = userHdAddress.ToAddress()

	validators := strings.Split(parts[1], ",")
	if len(validators)%2 != 0 {
		return fmt.Errorf("invalid format input[%s]", input)
	}
	for i := 0; i+1 < len(validators); i += 2 {
		validatorHdAddress, err := parseAddress(fmt.Sprintf("%s,%s", validators[i], validators[i+1]))
		if err != nil {
			return err
		}
		checkBalancePara.Validators = append(checkBalancePara.Validators, validatorHdAddress.ToAddress())
	}

	checkBalancePara.NetStake = big.NewInt(0)
	for i, v := range strings.Split(parts[2], ",") {
		val, ok := new(big.Int).SetString(v, 10)
		if !ok {
			return fmt.Errorf("invalid format input[%s]", input)
		}
		if i == 0 {
			checkBalancePara.NetStake = new(big.Int).Mul(checkBalancePara.NetStake, val)
		} else {
			checkBalancePara.NetStake = new(big.Int).Add(checkBalancePara.NetStake, val)
		}
	}
	c.rawAction.CheckBalancePara = checkBalancePara

	return nil
}

func (c *CheckBalanceParser) ParseAssertion(input string) error {
	return nil
}
