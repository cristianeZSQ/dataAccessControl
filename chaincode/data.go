/*
This is an Access Control Contract
Defines user and AccessRule


*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"time"
	"encoding/json"
	"strconv"
	"bytes"
	"strings"
)

type user struct {
	UserId int32 `json:"userId"`
	UserName string `json:"userName"`
	UserTel string `json:"userTel"`
	UserAddress string `json:"userAddress"`
}


type ACC struct {
	ObjectType string `json:"docType"`
	AccId int `json:"accId"`
	UserName string `json:"userName"`
	DataIndex string `json:"dataIndex"`
	AccessRule string `json:"accessRule"`
	Date time.Time `json:"date"`
}

type dataOrder struct {
	ObjectType string `json:"docType"`
	Owner string `json:"owner"`
	Requestor string `json:"requestor"`
	DataItem string `json:"dataItem"`
	Description string `json:"description"`
	Operation string `json:"operation"`
	State string `json:"state"`       // 0:denied 1:accessed  2:waiting
	Date time.Time `json:"date"`
}

//init acc table
//input: accId		userName   DataIndex	 rules	 date(using system currentTime)
func (t *SimpleChaincode) createACC(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	// 0				1			   2		        		3			4
	//accId		"userName"   "data_222"	 "{read:"professor Wang"}"	 "2018/12/22"
	if len(args) != 4{
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init access control contract")
	if len(args[0]) < 0{
		return  shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[2]) <= 5  {
		return shim.Error("the length of 2nd argument must be longer than 5")
	}
/*	result,err := regexp.Match(`^data_[0-9]*$`,[]byte(args[2]))
	if !result{
		return shim.Error(err.Error())
	}*/

	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	accId, err := strconv.Atoi(args[0])
	if err != nil{
		return shim.Error(err.Error())
	}
	userName := args[1]
	DataIndex := args[2]
	//init rules
	rules := args[3]
	date := time.Now()
	objectType := "acc"

	acc := &ACC{objectType,accId,userName,DataIndex,rules,date}
	accJSONasBytes, err := json.Marshal(acc)

	// === Save Acess control Rule table to state ===
	err = stub.PutState("ACC"+args[0], accJSONasBytes)
	if err != nil{
		return shim.Error(err.Error())
	}

	//construct key to query database

	//
	fmt.Println("-end init acc (success)")
	return shim.Success(nil)

}

//delete acc
//input: accId
func (t *SimpleChaincode) deleteACC(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	var jsonResp string
	var accJSON ACC
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accId := args[0]
	valAsbytes, err := stub.GetState(accId) //get the acc form chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + accId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + accId + "\"}"
		return shim.Error(jsonResp)
	}

	err=json.Unmarshal([]byte(valAsbytes),&accJSON)
	if err != nil{
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + accId + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(accId)
	if err != nil{
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

// add new rules to existed acc
func (t *SimpleChaincode) addRules(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accId := args[0]
	newRule := args[1]
	//check if acc already exists
	accAsBytes, err := stub.GetState(accId)
	if err != nil {
		return shim.Error("Failed to get acc: " + err.Error())
	} else if accAsBytes == nil {
		return shim.Error(accId + " does not exist ")
	}

	accNewRule := ACC{}
	err =json.Unmarshal(accAsBytes,&accNewRule)
	if err != nil {
		return shim.Error(err.Error())
	}

	accNewRule.AccessRule = accNewRule.AccessRule + ";" + newRule //change the rules

	accJSONasBytes, _ := json.Marshal(accNewRule)
	err = stub.PutState(accId,accJSONasBytes)
	if err!=nil{
		return shim.Error(err.Error())
	}

	fmt.Println("- end add New Rule to Acc (success)")
	return shim.Success(nil)

}

func (t *SimpleChaincode) changeRules(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accId := args[0]
	newRule := args[1]
	//check if acc already exists
	accAsBytes, err := stub.GetState(accId)
	if err != nil {
		return shim.Error("Failed to get acc: " + err.Error())
	} else if accAsBytes == nil {
		return shim.Error(accId + " does not exist ")
	}

	accNewRule := ACC{}
	err =json.Unmarshal(accAsBytes,&accNewRule)
	if err != nil {
		return shim.Error(err.Error())
	}

	accNewRule.AccessRule = newRule //change the rules

	accJSONasBytes, _ := json.Marshal(accNewRule)
	err = stub.PutState(accId,accJSONasBytes)
	if err!=nil{
		return shim.Error(err.Error())
	}

	fmt.Println("- end Change Rules to Acc (success)")
	return shim.Success(nil)
}

//if the storage place of personal data has been changed, use this method to update DataIndex
//input: accId, DataIndex  (ACC1,data_011)
func (t *SimpleChaincode) updateDataIndex(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accId := args[0]
	DataIndex := args[1]

	//check if acc already exists
	accAsBytes, err := stub.GetState(accId)
	if err != nil {
		return shim.Error("Failed to get acc: " + err.Error())
	} else if accAsBytes == nil {
		return shim.Error(accId + " does not exist ")
	}

	accNewIndex := ACC{}
	err =json.Unmarshal(accAsBytes,&accNewIndex)
	if err != nil {
		return shim.Error(err.Error())
	}

	accNewIndex.DataIndex = DataIndex //change DataIndex

	accJSONasBytes, _ := json.Marshal(accNewIndex)
	err = stub.PutState(accId,accJSONasBytes)
	if err!=nil{
		return shim.Error(err.Error())
	}

	fmt.Println("- end Change DataIndex to Acc (success)")

	return shim.Success(nil)
}

//input: userName "Bob"
func (t *SimpleChaincode) queryAccByOwner(stub shim.ChaincodeStubInterface,args []string) peer.Response{

	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	
	userName := strings.ToLower(args[0])
	
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"acc\",\"userName\":\"%s\"}}", userName)
	
	queryResults, err := getQueryResultForQueryString(stub,queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error)  {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}


// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface)(*bytes.Buffer,error){
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext(){
		queryResponse, err := resultsIterator.Next()
		if err != nil{
			return nil,err
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(",\"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer,nil
}

func (t *SimpleChaincode) queryAllACC(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	var err error

	startKey := "ACC1"
	endKey := "ACC999"

	resultsIterator, err := stub.GetStateByRange(startKey,endKey)
	if err != nil{
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil{
		return shim.Error(err.Error())
	}

	fmt.Printf("- getAllACC queryResult: \n%s \n",buffer.String())
	return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) getHistoryForAcc(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accId := args[0]


	fmt.Printf("- start getHistoryForAcc: acc_%s\n", accId)

	resultsIterator, err := stub.GetHistoryForKey(accId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

func (t *SimpleChaincode) getHistoryForRules(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	return shim.Success(nil)
}
