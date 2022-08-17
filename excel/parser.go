package excel

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/polynetwork/bridge-common/log"
	"strconv"
	"strings"
	"zion-test/base"
	"zion-test/zioncase"
)

type ComposeHandler interface {
	compose() error
}

func ParseExcel(excelPath string) {
	c := &CreateValidatorParam{}
	c.ConsensusAddress = zioncase.HDAddress{Index_1: 1, Index_2: 1}
	c.SignerAddress = zioncase.HDAddress{Index_1: 1, Index_2: 2}
	c.ProposalAddress = zioncase.HDAddress{Index_1: 1, Index_2: 3}
	c.Commission = 0
	c.InitStake = 0
	c.Desc = "validator 1"

	jsondata, _ := json.Marshal(c)
	fmt.Println("CreateValidatorParam=%s", string(jsondata))

	excel, err := excelize.OpenFile(excelPath)
	if err != nil {
		log.Fatal("open excel file failed", "err", err)
	}

	for i := 0; i < excel.SheetCount; i++ {
		for j, row := range excel.GetRows(excel.GetSheetName(i)) {
			if j == 0 {
				continue
			}
			if len(row) == 0 {
				continue
			}

			rowAction, err := createRowAction(row)
			_ = rowAction // todo
			if err != nil {
				log.Fatal("createRowAction failed", "err", err)
			}
		}
	}
}

func createRowAction(row []string) (action *zioncase.RawAction, err error) {
	if len(row) == 0 {
		return
	}
	action = &zioncase.RawAction{}

	caseNo := row[0]

	action.MethodName = row[1]
	action.RawInput = row[3]
	if row[4] == "1" {
		action.ShouldSucceed = true
	}
	if row[5] == "nil" {
		action.Result = nil
	}

	// parse sender
	senderSlice := strings.Split(row[6], ",")
	if len(senderSlice) != 2 {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	index1, err := strconv.ParseUint(senderSlice[0], 10, 64)
	if err != nil {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	index2, err := strconv.ParseUint(senderSlice[1], 10, 64)
	if err != nil {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	action.Sender = zioncase.HDAddress{Index_1: index1, Index_2: index2}

	// parse Options
	optionSlice := strings.Split(row[7], ",")
	if len(optionSlice) != 2 {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	block, err := strconv.ParseUint(senderSlice[0], 10, 64)
	if err != nil {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	shouldBefore, err := strconv.ParseUint(senderSlice[1], 10, 64)
	if err != nil {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	action.Options = zioncase.ActionBase{Block: block, ShouldBefore: shouldBefore}

	composeHandler := NewComposeHandler(action)
	_ = composeHandler
	action.Input = nil // todo

	switch action.MethodName {
	case base.MethodCreateValidator:
		_, _ = parseCreateValidatorParam(row[3]) // todo
	default:

	}

	//rowNo := row[0]

	//paras := row[2]
	//shouldSucceed := row[3]
	//result := row[4]
	//sender := row[5]
	//options := row[6]
	log.Info("row", "", row)
	return
}

func NewComposeHandler(rawAction *zioncase.RawAction) ComposeHandler {
	switch rawAction.MethodName {
	case base.MethodCreateValidator:
		return CreateValidatorComposer{rawAction: rawAction}

	default:
		log.Fatal("undefined method", "method", rawAction.MethodName)
		return nil

	}
}
