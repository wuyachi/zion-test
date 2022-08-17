package excel

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/polynetwork/bridge-common/log"
	"main/base"
	"main/common"
	"strconv"
	"strings"
)

type ComposeHandler interface {
	compose() error
}

func ParseExcel(excelPath string) {
	c := &CreateValidatorParam{}
	c.ConsensusAddress = common.HDAddress{Index_1: 1, Index_2: 1}
	c.SignerAddress = common.HDAddress{Index_1: 1, Index_2: 2}
	c.ProposalAddress = common.HDAddress{Index_1: 1, Index_2: 3}
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

func createRowAction(row []string) (action *common.RawAction, err error) {
	if len(row) == 0 {
		return
	}
	action = &common.RawAction{}

	caseNo := row[0]

	action.MethodName = row[1]
	//action.RawInput = row[3]
	if row[4] == "1" {
		action.ShouldSucceed = true
	}

	// parse sender
	senderSlice := strings.Split(row[6], ",")
	if len(senderSlice) != 2 {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	index1, err := strconv.ParseUint(senderSlice[0], 10, 32)
	if err != nil {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	index2, err := strconv.ParseUint(senderSlice[1], 10, 32)
	if err != nil {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	action.Sender = common.HDAddress{Index_1: uint32(index1), Index_2: uint32(index2)}

	// parse Options
	optionSlice := strings.Split(row[7], ",")
	if len(optionSlice) != 3 {
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
	epochId, err := strconv.ParseUint(senderSlice[2], 10, 64)
	if err != nil {
		log.Fatal("parse Sender failed", "Case No", caseNo, "Sender", row[6])
	}
	action.Block = block
	action.ShouldBefore = shouldBefore
	action.Epoch = epochId

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

func NewComposeHandler(rawAction *common.RawAction) ComposeHandler {
	switch rawAction.MethodName {
	case base.MethodCreateValidator:
		return CreateValidatorComposer{rawAction: rawAction}

	default:
		log.Fatal("undefined method", "method", rawAction.MethodName)
		return nil

	}
}
