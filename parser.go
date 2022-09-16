package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/polynetwork/bridge-common/log"
	"main/base"
	"strconv"
	"strings"
)

//excel column titles
const (
	_CaseNo        = "Case No"
	_CaseNote      = "caseNote"
	_Step          = "Step"
	_Note          = "Note"
	_ExpectDesc    = "ExpectDesc"
	_MethodName    = "MethodName"
	_ParamDesc     = "Param Desc"
	_Input         = "Input"
	_ShouldSucceed = "ShouldSucceed"
	_Assertion     = "Assertion"
	_Sender        = "Sender"
	_ActionBase    = "ActionBase"
	_ActionResult  = "ExeuteResult"
	_CaseResult    = "CaseResult"
)

type ParseHandler interface {
	ParseInput(input string) error
	ParseAssertion(input string) error
}

func NewParseHandler(rawAction *RawAction) (ParseHandler, error) {
	switch rawAction.MethodName {
	case "checkBalance":
		return &CheckBalanceParser{rawAction: rawAction}, nil
	case base.MethodCreateValidator:
		return &CreateValidatorParser{rawAction: rawAction}, nil
	case base.MethodWithdrawValidator:
		return &WithdrawValidatorParser{rawAction: rawAction}, nil
	case base.MethodUpdateValidator:
		return &UpdateValidatorParser{rawAction: rawAction}, nil
	case base.MethodCancelValidator:
		return &CancelValidatorParser{rawAction: rawAction}, nil
	case base.MethodUpdateCommission:
		return &UpdateCommissionParser{rawAction: rawAction}, nil
	case base.MethodWithdrawCommission:
		return &WithdrawCommissionParser{rawAction: rawAction}, nil
	case base.MethodStake:
		return &StakeParser{rawAction: rawAction}, nil
	case base.MethodUnStake:
		return &UnStakeParser{rawAction: rawAction}, nil
	case base.MethodWithdraw:
		return &WithdrawParser{rawAction: rawAction}, nil
	case base.MethodWithdrawStakeRewards:
		return &WithdrawStakeRewardsParser{rawAction: rawAction}, nil
	case base.MethodGetCurrentEpochInfo:
		return &GetCurrentEpochInfoParser{rawAction: rawAction}, nil
	case base.MethodGetAllValidators:
		return &GetAllValidatorsParser{rawAction: rawAction}, nil
	case base.MethodGetValidator:
		return &GetValidatorParser{rawAction: rawAction}, nil
	case base.MethodGetStakeInfo:
		return &GetStakeInfoParser{rawAction: rawAction}, nil
	case base.MethodGetStakeStartingInfo:
		return &GetStakeStartingInfoParser{rawAction: rawAction}, nil
	case base.MethodPropose:
		return &ProposeParser{rawAction: rawAction}, nil
	case base.MethodProposeConfig:
		return &ProposeConfigParser{rawAction: rawAction}, nil
	case base.MethodProposeCommunity:
		return &ProposeCommunityParser{rawAction: rawAction}, nil
	case base.MethodVoteProposal:
		return &VoteProposalParser{rawAction: rawAction}, nil
	case base.MethodGetProposal:
		return &GetProposalParser{rawAction: rawAction}, nil
	case base.MethodGetProposalList:
		return &GetProposalListParser{rawAction: rawAction}, nil
	case base.MethodGetConfigProposalList:
		return &GetConfigProposalListParser{rawAction: rawAction}, nil
	case base.MethodGetCommunityProposalList:
		return &GetCommunityProposalListParser{rawAction: rawAction}, nil
	default:
		err := fmt.Errorf("undefined method: %s", rawAction.MethodName)
		return nil, err
	}
}

func ParseExcel(excelPath string) (rawCases []*RawCase, err error) {
	excel, err := excelize.OpenFile(excelPath)
	if err != nil {
		log.Fatal("open excel file failed", "err", err)
	}
	for i := 1; i <= excel.SheetCount; i++ {
		var fieldsIndex map[string]int
		caseRows := make([][]string, 0)
		sheetName := excel.GetSheetName(i)
		if !strings.HasPrefix(sheetName, "case") {
			continue
		}
		rows := excel.GetRows(sheetName)
		for j, row := range rows {
			if j == 0 {
				fieldsIndex = getFieldsIndex(row)
				continue
			}
			if len(row) != 0 {
				caseRows = append(caseRows, row)
			}

			if len(rows) != j+1 {
				continue
			}

			// end of case
			rawCase, e := createRawCase(caseRows, fieldsIndex)
			if e != nil {
				return nil, e
			}
			rawCases = append(rawCases, rawCase)
			caseRows = make([][]string, 0)
		}
	}
	log.Info("Parsed excel", "sheet_count", excel.SheetCount, "case_count", len(rawCases))
	return
}

func createRawCase(rows [][]string, fieldsIndex map[string]int) (rawCase *RawCase, err error) {
	rawCase = &RawCase{Actions: []*RawAction{}}
	for i, row := range rows {
		if i == 0 {
			caseNo, err := strconv.ParseFloat(row[fieldsIndex[_CaseNo]], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid caseNo: %s", row[fieldsIndex[_CaseNo]])
			}
			rawCase.Index = int64(caseNo)
		}

		err = formatRow(row, fieldsIndex)
		if err != nil {
			err = fmt.Errorf("case format invalid, caseNo: %d, step: %d, err: %v", rawCase.Index, i+1, err)
			return
		}

		action, e := createRowAction(row, fieldsIndex)
		if e != nil {
			err = fmt.Errorf("createRowAction failed, caseNo: %d, step: %d, row: %s, err: %v", rawCase.Index, i+1, row, e)
			return
		}
		rawCase.Actions = append(rawCase.Actions, action)
	}
	return
}

func createRowAction(row []string, fieldsIndex map[string]int) (action *RawAction, err error) {

	action = new(RawAction)
	action.Row = row

	// MethodName
	action.MethodName = row[fieldsIndex[_MethodName]]

	// ShouldSucceed
	if row[fieldsIndex[_ShouldSucceed]] == "1" {
		action.ShouldSucceed = true
	}

	// Sender
	switch getMethodType(action.MethodName) {
	case TX:
		action.Sender, err = parseAddress(row[fieldsIndex[_Sender]])
		if err != nil {
			return
		}
	}

	// ActionBase
	action.Epoch, action.Block, action.ShouldBefore, err = parseActionBase(row[fieldsIndex[_ActionBase]])
	if err != nil {
		return
	}
	if action.MethodName == "checkBalance" {
		action.Block += 10 // delay 10 blocks
	}

	parseHandler, err := NewParseHandler(action)
	if err != nil {
		err = fmt.Errorf("new parseHandler failed, err: %s", err)
		return
	}

	// Input
	err = parseHandler.ParseInput(row[fieldsIndex[_Input]])
	if err != nil {
		err = fmt.Errorf("parse Input failed, err: %s", err)
		return nil, err
	}

	// Assertion
	err = parseHandler.ParseAssertion(row[fieldsIndex[_Assertion]])
	if err != nil {
		err = fmt.Errorf("parse Input failed, err: %s", err)
		return nil, err
	}

	return
}

func getFieldsIndex(fields []string) map[string]int {
	fieldsIndex := make(map[string]int, 0)
	for i, field := range fields {
		fieldsIndex[field] = i
	}
	return fieldsIndex
}

func formatRow(row []string, fieldsIndex map[string]int) error {
	for i := 0; i < len(row); i++ {
		row[i] = strings.Replace(row[i], "[", "", -1)
		row[i] = strings.Replace(row[i], "]", "", -1)
		row[i] = strings.Replace(row[i], "]", "\"", -1)
	}
	return checkRow(row, fieldsIndex)
}

func checkRow(row []string, fieldsIndex map[string]int) (err error) {
	// check MethodName
	methodName := row[fieldsIndex[_MethodName]]
	if len(methodName) == 0 {
		err = fmt.Errorf("invalid format [MethodName]: %s", methodName)
		return
	}

	// check ActionBase
	actionBase := row[fieldsIndex[_ActionBase]]
	parts := strings.Split(actionBase, ",")
	if len(parts) != 3 {
		err = fmt.Errorf("invalid format [ActionBase]: %s", actionBase)
		return
	}
	// check Assertion
	assertion := row[fieldsIndex[_Assertion]]
	if assertion != "nil" {
		parts = strings.Split(assertion, ";")
		if len(parts) < 3 {
			err = fmt.Errorf("invalid format [Assertion]: %s", assertion)
			return
		}
	}
	// check Sender
	sender := row[fieldsIndex[_Sender]]
	if sender != "nil" {
		parts = strings.Split(sender, ",")
		if len(parts) != 2 {
			err = fmt.Errorf("invalid format [Sender]: %s", sender)
			return
		}
	}
	return
}

func parseAddress(input string) (address HDAddress, err error) {
	parts := strings.Split(input, ",")
	index1, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		err = fmt.Errorf("parse address failed, input: %s, err: %v", input, err)
		return
	}
	index2, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		err = fmt.Errorf("parse address failed, input: %s, err: %v", input, err)
		return
	}

	address.Index_1 = uint32(index1)
	address.Index_2 = uint32(index2)
	return
}

func parseActionBase(input string) (epoch, block, shouldBefore uint64, err error) {
	parts := strings.Split(input, ",")
	epoch, err = strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		err = fmt.Errorf("parse actionBase failed, input: %s", input)
		return
	}
	block, err = strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		err = fmt.Errorf("parse actionBase failed, input: %s", input)
		return
	}
	shouldBefore, err = strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		err = fmt.Errorf("parse actionBase failed, input: %s", input)
		return
	}
	return
}

func formatAssertType(tag string) (assertType AssertType, err error) {
	switch tag {
	case "contain":
		assertType = Assert_Element_Contain
	case "notContain":
		assertType = Assert_Element_Not_Contain
	case "equal":
		assertType = Assert_Element_Equal
	case "notEqual":
		assertType = Assert_Element_Not_Equal
	default:
		err = fmt.Errorf("undefined assert type: %s", tag)
	}
	return
}
