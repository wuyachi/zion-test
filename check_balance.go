package main

import (
	"fmt"
	"strings"
)

type CheckBalanceParser struct {
	rawAction *RawAction
}

func (c *CheckBalanceParser) ParseInput(input string) error {
	parts := strings.Split(input, ";")
	if len(parts) == 0 {
		return fmt.Errorf("invalid format input[%s]", input)
	}

	for _, part := range parts {
		hdAddress, err := parseAddress(part)
		if err != nil {
			return err
		}
		address := hdAddress.ToAddress()
		c.rawAction.BalanceAddresses = append(c.rawAction.BalanceAddresses, address)
	}

	return nil
}

func (c *CheckBalanceParser) ParseAssertion(input string) error {
	return nil
}
