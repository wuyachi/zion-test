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
	_ExeuteResult  = "ExeuteResult"
)

type ParseHandler interface {
	ParseInput(input string) (Param, error)
	ParseAssertion(input string) ([]Assertion, error)
}

func NewParseHandler(rawAction *RawAction) (ParseHandler, error) {
	switch rawAction.MethodName {
	case base.MethodCreateValidator:
		return &CreateValidatorParser{rawAction: rawAction}, nil
	case base.MethodGetCurrentEpochInfo:
		return &GetCurrentEpochInfoParser{rawAction: rawAction}, nil
	case base.MethodGetAllValidators:
		return &GetAllValidatorsParser{rawAction: rawAction}, nil
	case base.MethodStake:
		return &StakeParser{rawAction: rawAction}, nil
	case base.MethodUnStake:
		return &UnStakeParser{rawAction: rawAction}, nil
	case base.MethodWithdraw:
		return &WithdrawParser{rawAction: rawAction}, nil
	case base.MethodWithdrawStakeRewards:
		return &WithdrawStakeRewardsParser{rawAction: rawAction}, nil
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
		rows := excel.GetRows(excel.GetSheetName(i))
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
				log.Error("createRawCase failed", "err", e)
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
			rawCase.Index, err = strconv.ParseInt(row[fieldsIndex[_CaseNo]], 10, 64)
			if err != nil {
				err = fmt.Errorf("invalid caseNo: %s", row[fieldsIndex[_CaseNo]])
				return
			}
		}

		err = formatRow(row, fieldsIndex)
		if err != nil {
			err = fmt.Errorf("case format invalid, caseNo: %d, step: %d, err: %v", rawCase.Index, i+1, err)
			return
		}

		action, e := createRowAction(row, fieldsIndex)
		if e != nil {
			err = fmt.Errorf("createRowAction failed, caseNo: %d, row: %s, err: %v", rawCase.Index, row, e)
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
	if !ReadOnly(action.MethodName) {
		action.Sender, err = parseAddress(row[fieldsIndex[_Sender]])
		if err != nil {
			err = fmt.Errorf("parse Sender failed, Sender: %s", row[fieldsIndex[_Sender]])
			return
		}
	}

	// ActionBase
	action.Epoch, action.Block, action.ShouldBefore, err = parseActionBase(row[fieldsIndex[_ActionBase]])
	if err != nil {
		return
	}

	parseHandler, err := NewParseHandler(action)
	if err != nil {
		err = fmt.Errorf("new parseHandler failed, err: %s", err)
		return
	}

	// Input
	action.Input, err = parseHandler.ParseInput(row[fieldsIndex[_Input]])
	if err != nil {
		err = fmt.Errorf("parse Input failed, err: %s", err)
		return nil, err
	}

	// Assertion
	action.Assertions, err = parseHandler.ParseAssertion(row[fieldsIndex[_Assertion]])
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
		parts = strings.Split(assertion, ",")
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
		err = fmt.Errorf("parse address failed, input: %s", input)
		return
	}
	index2, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		err = fmt.Errorf("parse address failed, input: %s", input)
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
