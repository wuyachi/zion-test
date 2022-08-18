package main

import (
	"fmt"
	"github.com/devfans/zion-sdk/contracts/native/go_abi/proposal_manager_abi"
	"github.com/devfans/zion-sdk/contracts/native/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"math/big"
	"strings"
)

func init() {
	ab, err := abi.JSON(strings.NewReader(proposal_manager_abi.IProposalManagerMetaData.ABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load proposal manager abi json string: [%v]", err))
	}
	ProposalManagerABI = &ab
}

var (
	ProposalManagerABI *abi.ABI
)

type ProposeParam struct {
	Content []byte
}

func (m *ProposeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ProposalManagerABI, MethodPropose, m)
}

type ProposeConfigParam struct {
	Content []byte
}

func (m *ProposeConfigParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ProposalManagerABI, MethodProposeConfig, m)
}

type ProposeCommunityParam struct {
	Content []byte
}

func (m *ProposeCommunityParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ProposalManagerABI, MethodProposeCommunity, m)
}

type VoteProposalParam struct {
	ID *big.Int
}

func (m *VoteProposalParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ProposalManagerABI, MethodVoteProposal, m)
}

type GetProposalParam struct {
	ID *big.Int
}

func (m *GetProposalParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ProposalManagerABI, MethodGetProposal, m)
}

type GetProposalListParam struct{}

func (m *GetProposalListParam) Encode() ([]byte, error) {
	return utils.PackMethod(ProposalManagerABI, MethodGetProposalList)
}

type GetConfigProposalListParam struct{}

func (m *GetConfigProposalListParam) Encode() ([]byte, error) {
	return utils.PackMethod(ProposalManagerABI, MethodGetConfigProposalList)
}

type GetCommunityProposalListParam struct{}

func (m *GetCommunityProposalListParam) Encode() ([]byte, error) {
	return utils.PackMethod(ProposalManagerABI, MethodGetCommunityProposalList)
}
