package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"strings"
)


//user2 ask for user1's personal infomation
//( u2 , u1 , DataItem,Description,Operation)
func (t *SimpleChaincode) createDataOrder(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	// 0				1			   2		      	3			4
	//OwnerName		Requestor     DataItem	       Description	 Operation
	if len(args) != 5{
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start to create an order to request someone's data")
	if len(args[0]) < 0{
		return  shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) <= 0 {
		return shim.Error("2st argument must be a non-empty string")
	}

	if len(args[2]) <= 0  {
		return shim.Error("3st argument must be a non-empty string")
	}

	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	if len(args[4]) > 0 {
		if strings.ToLower(args[4]) != "read" && strings.ToLower(args[4]) != "write"{
			return shim.Error("operation must be 'read' OR 'write'")
		}
	}else{
		return shim.Error("4th argument must be a non-empty string")
	}

	objectType := "dataOrder"

	dataOrder := &dataOrder{objectType,args[0],args[1],args[2],args[3],args[4],"waiting",time.Now()}
	dataOrderAsBytes, err := json.Marshal(dataOrder)

	// === Save table to state ===
	err = stub.PutState("Order_"+args[1]+args[2], dataOrderAsBytes)
	if err != nil{
		return shim.Error(err.Error())
	}

	/*
	//construct key to query database
	indexName := "dataItemIndex"
	dataItemIndexKey, err := stub.CreateCompositeKey(indexName,[]string{dataOrder.DataItem})

	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value

	stub.PutState(dataItemIndexKey, dataOrderAsBytes)
	*/

	fmt.Println("-end create dataOrder (success)")
	return shim.Success(nil)

}

//input:requestor, dataItem
func (t *SimpleChaincode)submitOrder(stub shim.ChaincodeStubInterface, args []string) peer.Response{
	if len(args) != 2{
		shim.Error("Incorrect number of arguments. Expecting requestor and dataItem")
	}

	orderAsBytes, err := stub.GetState("Order_"+args[0]+args[1])
	if err != nil {
		return shim.Error("Failed to get order: " + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error(args[0] + " does not exist ")
	}

	Order := dataOrder{}
	err = json.Unmarshal(orderAsBytes,&Order)
	if err != nil{
		return shim.Error(err.Error())
	}

	requestor := args[0]
	dataItem := args[1]
	newState := "state wrong"

	fmt.Println("- start changeStateBasedOndataItem ", dataItem)

	var result bool
	result, err = checkRequestRule(requestor,dataItem,stub)

	if result {
		newState = "accessed"
	}else{
		newState = "denied"
	}

	Order.State = newState

	inputAsBytes, _ := json.Marshal(Order)
	err = stub.PutState("Order_" + args[0]+args[1],inputAsBytes)
	if err!=nil{
		return shim.Error(err.Error())
	}

	fmt.Println("- end submit Order (success)")
	return shim.Success(nil)

}

func checkRequestRule(requestor string,dataItem string, stub shim.ChaincodeStubInterface) (bool,error){
	var state bool

/*	input := struct {
		Requestor string `json:"requestor"`
		DataItem string `json:"dataItem"`
	}{}

	indexName := "dataItemIndex"
	orderKey, err := stub.CreateCompositeKey(indexName,[]string{input.DataItem})
	if err != nil{
		fmt.Println(err.Error())
		return false,err
	}*/

	orderBytes, err := stub.GetState("Order_"+requestor+dataItem)
	if len(orderBytes) == 0{
		fmt.Println("Could not find the dataOrder")
		return false,err
	}
	//解析dataOrder
	dataOrder := dataOrder{}
	err = json.Unmarshal(orderBytes,&dataOrder)
	if err != nil{
		fmt.Println("dataOrder 反序列化 失败")
		return false,err
	}

	//利用从数据库中得到的结果构建一个rule数据结构
	newRule := dataOrder.Operation + ":" + dataOrder.Requestor
	owner := dataOrder.Owner

	//查找owner的ACC中是否存在这样的rule
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"acc\",\"userName\":\"%s\",\"dataIndex\":\"%s\"}}", owner,dataItem)

	queryResults, err := getQueryResultForQueryString(stub,queryString)
	if err != nil {
		fmt.Printf("Could not find a acc by input userName:%s and dataIndex:%s",owner,dataItem)
		return false,err
	}

	//这里可以对queryResult进行解析
	ownerACC := ACC{}
	err = json.Unmarshal(queryResults,&ownerACC)
	if err != nil{
		return false, err
	}

	state = strings.Contains(ownerACC.AccessRule,newRule)

	return state,nil
}

func (t *SimpleChaincode)getAllOrder(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	var err error

	startKey := "Order_a"
	endKey := "Order_z"

	resultsIterator, err := stub.GetStateByRange(startKey,endKey)
	if err != nil{
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil{
		return shim.Error(err.Error())
	}

	fmt.Printf("- getAllOrder queryResult: \n%s \n",buffer.String())
	return shim.Success(buffer.Bytes())
}


func listOrderByRequestor(stub shim.ChaincodeStubInterface) peer.Response{


	return shim.Success(nil)
}

func listOrderByOwner() peer.Response{
	return shim.Success(nil)
}
