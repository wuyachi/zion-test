package main

import (
	"fmt"
	"github.com/devfans/zion-sdk/contracts/native/go_abi/node_manager_abi"
	"github.com/devfans/zion-sdk/contracts/native/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

func init() {
	ab, err := abi.JSON(strings.NewReader(node_manager_abi.INodeManagerMetaData.ABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load node manager abi json string: [%v]", err))
	}
	NodeManagerABI = &ab
}

var (
	NodeManagerABI *abi.ABI
)

type CreateValidatorParam struct {
	ConsensusAddress common.Address
	SignerAddress    common.Address
	ProposalAddress  common.Address
	Commission       *big.Int
	InitStake        *big.Int
	Desc             string
}

func (m *CreateValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodCreateValidator, m)
}

type UpdateValidatorParam struct {
	ConsensusAddress common.Address
	ProposalAddress  common.Address
	Desc             string
}

func (m *UpdateValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodUpdateValidator, m)
}

type UpdateCommissionParam struct {
	ConsensusAddress common.Address
	Commission       *big.Int
}

func (m *UpdateCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodUpdateCommission, m)
}

type StakeParam struct {
	ConsensusAddress common.Address
	Amount           *big.Int
}

func (m *StakeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodStake, m)
}

type UnStakeParam struct {
	ConsensusAddress common.Address
	Amount           *big.Int
}

func (m *UnStakeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodUnStake, m)
}

type CancelValidatorParam struct {
	ConsensusAddress common.Address
}

func (m *CancelValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodCancelValidator, m)
}

type WithdrawValidatorParam struct {
	ConsensusAddress common.Address
}

func (m *WithdrawValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodWithdrawValidator, m)
}

type WithdrawStakeRewardsParam struct {
	ConsensusAddress common.Address
}

func (m *WithdrawStakeRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodWithdrawStakeRewards, m)
}

type WithdrawCommissionParam struct {
	ConsensusAddress common.Address
}

func (m *WithdrawCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodWithdrawCommission, m)
}

type ChangeEpochParam struct{}

func (m *ChangeEpochParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodChangeEpoch)
}

type WithdrawParam struct{}

func (m *WithdrawParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodWithdraw)
}

type EndBlockParam struct{}

func (m *EndBlockParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodEndBlock)
}

type GetGlobalConfigParam struct{}

func (m *GetGlobalConfigParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodGetGlobalConfig)
}

type GetCommunityInfoParam struct{}

func (m *GetCommunityInfoParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodGetCommunityInfo)
}

type GetCurrentEpochInfoParam struct{}

func (m *GetCurrentEpochInfoParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodGetCurrentEpochInfo)
}

type GetEpochInfoParam struct {
	ID *big.Int
}

func (m *GetEpochInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetEpochInfo, m)
}

type GetAllValidatorsParam struct{}

func (m *GetAllValidatorsParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodGetAllValidators)
}

type GetValidatorParam struct {
	ConsensusAddress common.Address
}

func (m *GetValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetValidator, m)
}

type GetStakeInfoParam struct {
	ConsensusAddress common.Address
	StakeAddress     common.Address
}

func (m *GetStakeInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetStakeInfo, m)
}

type GetUnlockingInfoParam struct {
	StakeAddress common.Address
}

func (m *GetUnlockingInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetUnlockingInfo, m)
}

type GetStakeStartingInfoParam struct {
	ConsensusAddress common.Address
	StakeAddress     common.Address
}

func (m *GetStakeStartingInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetStakeStartingInfo, m)
}

type GetAccumulatedCommissionParam struct {
	ConsensusAddress common.Address
}

func (m *GetAccumulatedCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetAccumulatedCommission, m)
}

type GetValidatorSnapshotRewardsParam struct {
	ConsensusAddress common.Address
	Period           uint64
}

func (m *GetValidatorSnapshotRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetValidatorSnapshotRewards, m)
}

type GetValidatorAccumulatedRewardsParam struct {
	ConsensusAddress common.Address
}

func (m *GetValidatorAccumulatedRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetValidatorAccumulatedRewards, m)
}

type GetValidatorOutstandingRewardsParam struct {
	ConsensusAddress common.Address
}

func (m *GetValidatorOutstandingRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(NodeManagerABI, MethodGetValidatorOutstandingRewards, m)
}

type GetTotalPoolParam struct{}

func (m *GetTotalPoolParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodGetTotalPool)
}

type GetOutstandingRewardsParam struct{}

func (m *GetOutstandingRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethod(NodeManagerABI, MethodGetOutstandingRewards)
}
