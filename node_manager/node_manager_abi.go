package node_manager

import (
	"fmt"
	"github.com/devfans/zion-sdk/contracts/native/go_abi/node_manager_abi"
	"github.com/devfans/zion-sdk/contracts/native/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"main/base"
	"math/big"
	"strings"
)

func init() {
	ab, err := abi.JSON(strings.NewReader(node_manager_abi.INodeManagerMetaData.ABI))
	if err != nil {
		panic(fmt.Sprintf("failed to load node manager abi json string: [%v]", err))
	}
	ABI = &ab
}

var (
	ABI *abi.ABI
)

type CreateValidatorParam struct {
	ConsensusAddress common.Address
	SignerAddress    common.Address
	ProposalAddress  common.Address
	Commission       *big.Int
	Desc             string
}

func (m *CreateValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodCreateValidator, m)
}

type UpdateValidatorParam struct {
	ConsensusAddress common.Address
	SignerAddress    common.Address
	ProposalAddress  common.Address
	Desc             string
}

func (m *UpdateValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodUpdateValidator, m)
}

type UpdateCommissionParam struct {
	ConsensusAddress common.Address
	Commission       *big.Int
}

func (m *UpdateCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodUpdateCommission, m)
}

type StakeParam struct {
	ConsensusAddress common.Address
}

func (m *StakeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodStake, m)
}

type UnStakeParam struct {
	ConsensusAddress common.Address
	Amount           *big.Int
}

func (m *UnStakeParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodUnStake, m)
}

type CancelValidatorParam struct {
	ConsensusAddress common.Address
}

func (m *CancelValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodCancelValidator, m)
}

type WithdrawValidatorParam struct {
	ConsensusAddress common.Address
}

func (m *WithdrawValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodWithdrawValidator, m)
}

type WithdrawStakeRewardsParam struct {
	ConsensusAddress common.Address
}

func (m *WithdrawStakeRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodWithdrawStakeRewards, m)
}

type WithdrawCommissionParam struct {
	ConsensusAddress common.Address
}

func (m *WithdrawCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodWithdrawCommission, m)
}

type ChangeEpochParam struct{}

func (m *ChangeEpochParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodChangeEpoch)
}

type WithdrawParam struct{}

func (m *WithdrawParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodWithdraw)
}

type EndBlockParam struct{}

func (m *EndBlockParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodEndBlock)
}

type GetGlobalConfigParam struct{}

func (m *GetGlobalConfigParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodGetGlobalConfig)
}

type GetCommunityInfoParam struct{}

func (m *GetCommunityInfoParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodGetCommunityInfo)
}

type GetCurrentEpochInfoParam struct{}

func (m *GetCurrentEpochInfoParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodGetCurrentEpochInfo)
}

type GetEpochInfoParam struct {
	ID *big.Int
}

func (m *GetEpochInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetEpochInfo, m)
}

type GetAllValidatorsParam struct{}

func (m *GetAllValidatorsParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodGetAllValidators)
}

type GetValidatorParam struct {
	ConsensusAddress common.Address
}

func (m *GetValidatorParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetValidator, m)
}

type GetStakeInfoParam struct {
	ConsensusAddress common.Address
	StakeAddress     common.Address
}

func (m *GetStakeInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetStakeInfo, m)
}

type GetUnlockingInfoParam struct {
	StakeAddress common.Address
}

func (m *GetUnlockingInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetUnlockingInfo, m)
}

type GetStakeStartingInfoParam struct {
	ConsensusAddress common.Address
	StakeAddress     common.Address
}

func (m *GetStakeStartingInfoParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetStakeStartingInfo, m)
}

type GetAccumulatedCommissionParam struct {
	ConsensusAddress common.Address
}

func (m *GetAccumulatedCommissionParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetAccumulatedCommission, m)
}

type GetValidatorSnapshotRewardsParam struct {
	ConsensusAddress common.Address
	Period           uint64
}

func (m *GetValidatorSnapshotRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetValidatorSnapshotRewards, m)
}

type GetValidatorAccumulatedRewardsParam struct {
	ConsensusAddress common.Address
}

func (m *GetValidatorAccumulatedRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetValidatorAccumulatedRewards, m)
}

type GetValidatorOutstandingRewardsParam struct {
	ConsensusAddress common.Address
}

func (m *GetValidatorOutstandingRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetValidatorOutstandingRewards, m)
}

type GetTotalPoolParam struct{}

func (m *GetTotalPoolParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodGetTotalPool)
}

type GetOutstandingRewardsParam struct{}

func (m *GetOutstandingRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethod(ABI, base.MethodGetOutstandingRewards)
}

type GetStakeRewardsParam struct {
	ConsensusAddress common.Address
	StakeAddress     common.Address
}

func (m *GetStakeRewardsParam) Encode() ([]byte, error) {
	return utils.PackMethodWithStruct(ABI, base.MethodGetStakeRewards, m)
}
