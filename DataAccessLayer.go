package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func createDatabase(stub shim.ChaincodeStubInterface, args []string) (bool, error) {
	var err error
	//Create table "ContractDetails"
	err = stub.CreateTable("contractDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ContractId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "ContractObject", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return false, errors.New("Failed creating ContractDetails table")
	}

	err = stub.CreateTable("attachmentDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "contractId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "attachmentName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "documentBlob", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return false, errors.New("Failed creating attachmentDetails table.")
	}

	err = stub.CreateTable("userDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "userId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "contractIdList", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return false, errors.New("Failed creating userDetails table.")
	}

	return true, nil

}

func insertUserBlankRecord(stub shim.ChaincodeStubInterface, userId string) (bool, error) {
	var blankList []string
	var ok bool

	ok = insertUserContractList(stub, userId, blankList)
	if !ok {
		return ok, errors.New("Error in creating User")
	}
	return true, nil

}

func insertContractDetails(stub shim.ChaincodeStubInterface, contractDetails contract) (bool, error) {
	var err error
	var ok bool
	jsonAsBytes, _ := json.Marshal(contractDetails)
	ok, err = stub.InsertRow("contractDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: contractDetails.ContractId}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: jsonAsBytes}},
		},
	})
	return ok, err
}

func insertAttachmentDetails(stub shim.ChaincodeStubInterface, contractID string, attachmentName string, documentBlob string) (bool, error) {
	var err error
	var ok bool

	ok, err = stub.InsertRow("attachmentDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: contractID}},
			&shim.Column{Value: &shim.Column_String_{String_: attachmentName}},
			&shim.Column{Value: &shim.Column_String_{String_: documentBlob}},
		},
	})
	return ok, err
}

func getContractDetails(stub shim.ChaincodeStubInterface, contractId string) (contract, error) {

	var columns []shim.Column
	var contractList contract

	col1 := shim.Column{Value: &shim.Column_String_{String_: contractId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("contractDetails", columns)
	if err != nil {
		return contractList, errors.New("Failed to query table contractDetails")
	}
	contractAsBytes := row.Columns[1].GetBytes()
	json.Unmarshal(contractAsBytes, &contractList)

	return contractList, nil

}

func getAttachmentDetails(stub shim.ChaincodeStubInterface, contractId string, attachmentName string) ([]byte, error) {
	var documentBlob string
	var err error
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: contractId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: attachmentName}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	row, err := stub.GetRow("attachmentDetails", columns)
	if err != nil {
		return nil, errors.New("Failed to query table contractDetails")
	}
	documentBlob = row.Columns[2].GetString_()
	documentBlobAsBytes, _ := json.Marshal(documentBlob)
	return documentBlobAsBytes, nil
}

func getUserContractList(stub shim.ChaincodeStubInterface, userId string) ([]string, bool) {
	var columns []shim.Column
	var contractList []string

	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("userDetails", columns)
	if err != nil {
		return nil, false
	}

	contractListAsBytes := row.Columns[1].GetBytes()
	json.Unmarshal(contractListAsBytes, &contractList)

	return contractList, true
}

func updateUserContractList(stub shim.ChaincodeStubInterface, userId string, contractList []string) bool {
	JsonAsBytes, _ := json.Marshal(contractList)

	ok, err := stub.ReplaceRow("userDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: userId}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})

	if !ok && err == nil {
		return false
	}

	return true
}

func insertUserContractList(stub shim.ChaincodeStubInterface, userId string, contractList []string) bool {
	JsonAsBytes, _ := json.Marshal(contractList)

	ok, err := stub.InsertRow("userDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: userId}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})

	if !ok && err == nil {
		return false
	}
	return true
}

func updateContractListByContractID(stub shim.ChaincodeStubInterface, contractId string, contractList contract) bool {
	JsonAsBytes, _ := json.Marshal(contractList)

	ok, err := stub.ReplaceRow("contractDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: contractId}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})

	if !ok && err == nil {
		return false
	}

	return true
}

/*func GetUserSpecificContractList(stub shim.ChaincodeStubInterface, UserId string) ([]string, error) {
	var columns []shim.Column
	var ContractList []string
	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	columns = append(columns, col1)
	row, err := stub.GetRow("UserDetails", columns)
	if err != nil {
		return ContractList, errors.New("Failed to query table BuyerDetails")
	}
	json.Unmarshal(row.Columns[5].GetBytes(), &ContractList)
	return ContractList, nil
}*/
/*
func updateContractDetails(stub shim.ChaincodeStubInterface, contractDetails contract) (bool, error) {
	ok, err := stub.ReplaceRow("ContractDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.ContractId}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.OrderId}},
			&shim.Column{Value: &shim.Column_Bool{Bool: ContractDetails.PaymentCommitment}},
			&shim.Column{Value: &shim.Column_Bool{Bool: ContractDetails.PaymentConfirmation}},
			&shim.Column{Value: &shim.Column_Bool{Bool: ContractDetails.InformationCounterparty}},
			&shim.Column{Value: &shim.Column_Bool{Bool: ContractDetails.ForfeitingInvoice}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.ContractCreateDate}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.PaymentDueDate}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.InvoiceStatus}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.PaymentStatus}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.ContractStatus}},
			&shim.Column{Value: &shim.Column_String_{String_: ContractDetails.DeliveryStatus}},
		},
	})
	if !ok && err == nil {
		return false, errors.New("Error in updating Seller record.")
	}
	return true, nil
}
*/
