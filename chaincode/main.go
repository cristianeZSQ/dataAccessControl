

//====CHAINCODE EXECUTION SAMPLES (CLI) ==================
//
//==== Invoke function ====
//peer chaincode invoke -n mycc -c '{"Args":["createACC","1","bob","data_001","read:zhou"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["addRules","ACC1","read:alice"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["queryAllACC"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["getHistoryForAcc","ACC1"]}' -C myc
//peer chaincode invoke -n mycc -c '{"Args":["createDataOrder","bob","alice","data_001","bob personal information","read"]}' -C myc
//peer chaincode invoke -n mycc -C myc -c '{"Args":["submitOrder","alice","data_001"]}'
//peer chaincode invoke -n mycc -C myc -c '{"Args":["getAllOrder"]}'

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
)


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createACC"{
		return t.createACC(stub, args)
	}else if function == "deleteACC"{
		return t.deleteACC(stub,args)
	}else if function == "addRules"{
		return t.addRules(stub, args)
	}else if function == "changeRules"{
		return t.changeRules(stub, args)
	} else if function == "updateDataPlace"{
		return t.updateDataIndex(stub, args)
	} else if function == "queryAccByOwner"{
		return t.queryAccByOwner(stub, args)
	} else if function == "queryAllACC"{
		return t.queryAllACC(stub, args)
	} else if function == "createDataOrder"{
		return t.createDataOrder(stub,args)
	} else if function == "submitOrder"{
		return t.submitOrder(stub,args)
	}else if function == "getAllOrder" {
		return t.getAllOrder(stub,args)
	}else if function == "getHistoryForAcc"{
		return t.getHistoryForAcc(stub, args)
	} else if function == "getHistoryForRules"{
		return t.getHistoryForRules(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}
