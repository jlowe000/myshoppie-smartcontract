
package main

import (
	"fmt"
//	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//==============================================================================================================================
//	Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}		

//==============================================================================================================================
//	Account - Defines the structure for a MyLittleShopper object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON owner -> Struct Owner
//==============================================================================================================================
type MyLittleShopper struct{
	ShoppieNo string `json:"shoppieno"`	
	Owner string `json:"owner"`	
	Name string `json:"name"`
	Transaction string `json:"transaction"`
}

var shoppieIndexStr = "_shoppieindex"	  // Define an index varibale to track all the accounts stored in the world state

// ============================================================================================================================
//  Main - main - Starts up the chaincode
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init Function - Called when the user deploys the chaincode
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var errorMessage string

	// Init Some Little Shoppers
	
	var shop1 [4]string
	shop1[0] = "1" // "Unique Id"
	shop1[1] = "Coles" // "Coles owns them initially"
	shop1[2] = "Lipton Black Tea" // Name
	shop1[3] = "0.00" // What it is currently worth (to buy / transfer)

	errorMessage = t.init_new_shoppie(stub, shop1[:])
	if(errorMessage != "") {
		return shim.Error(errorMessage)
	}

	var shop2 [4]string
	shop2[0] = "2" // "Unique Id"
	shop2[1] = "Coles" // "Coles owns them initially"
	shop2[2] = "Huggies Nappies" // Name
	shop2[3] = "0.00" // What it is currently worth (to buy / transfer)

	errorMessage = t.init_new_shoppie(stub, shop2[:])
	if(errorMessage != "") {
		return shim.Error(errorMessage)
	}

	var shop3 [4]string
	shop3[0] = "3" // "Unique Id"
	shop3[1] = "Coles" // "Coles owns them initially"
	shop3[2] = "Weet-Bix Family Pack" // Name
	shop3[3] = "0.00" // What it is currently worth (to buy / transfer)

	errorMessage = t.init_new_shoppie(stub, shop3[:])
	if(errorMessage != "") {
		return shim.Error(errorMessage)
	}

	var shop4 [4]string
	shop4[0] = "4" // "Unique Id"
	shop4[1] = "Coles" // "Coles owns them initially"
	shop4[2] = "Sun Bites Snack Crackers" // Name
	shop4[3] = "0.00" // What it is currently worth (to buy / transfer)

	errorMessage = t.init_new_shoppie(stub, shop4[:])
	if(errorMessage != "") {
		return shim.Error(errorMessage)
	}

	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		    initial arguments passed to other things for use in the called function.
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "read" {
		return t.read(stub,args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "write" {
		return t.write(stub, args)
	}

	fmt.Println("query did not find func: " + function)						//error
	return shim.Error("Received unknown function invocation: " + function)
}

// ============================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetStringArgs()

	return t.read(stub,args)
}

// ============================================================================================================================
// Delete - remove a key/value pair from the world state
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return t.delete(stub,args)
}

// ============================================================================================================================
// Write - directly write a variable into chaincode world state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string)  pb.Response {

	// "write"
	return t.write(stub,args)
}

// ============================================================================================================================
// Read - read a variable from chaincode world state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)	
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	}
	
	return shim.Success(valAsbytes)												
}

// ============================================================================================================================
// Delete - remove a key/value pair from the world state
// ============================================================================================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	//get the shoppie index
	shoppieAsBytes, err := stub.GetState(shoppieIndexStr)
	if err != nil {
		return shim.Error("Failed to get shoppie index")
	}
	var shoppieIndex []string
	json.Unmarshal(shoppieAsBytes, &shoppieIndex)						
	
	//remove shoppie from world (index)
	for i,val := range shoppieIndex{
		if val == name{															//find the correct shoppie
			shoppieIndex = append(shoppieIndex[:i], shoppieIndex[i+1:]...)			//remove it
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(shoppieIndex)									//save the new index
	err = stub.PutState(shoppieIndexStr, jsonAsBytes)
	return shim.Success(nil)
}

// ============================================================================================================================
// Write - directly write a variable into chaincode world state
// ============================================================================================================================
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string)  pb.Response {
	var name, value string 
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]														
	value = args[1]
	err = stub.PutState(name, []byte(value))					
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// ============================================================================================================================
// Init shoppie - create a new shoppie, store into chaincode world state, and then append the shoppie index
// ============================================================================================================================
func (t *SimpleChaincode) init_new_shoppie(stub shim.ChaincodeStubInterface, args []string) (string) {
	var err error

	shoppieNo := args[0]

	//check if shoppie already exists
	// shoppieAsBytes, err := stub.GetState(shoppieNo)
	_, err = stub.GetState(shoppieNo)
	if err != nil {
		return "Failed to get shoppieNo number"
	}
	str := `{ "shoppieno": "` + args[0] + `", "owner": "` + args[1] + `", "name": "` + args[2] + `", "transaction": "` + args[3] + `" }`
	err = stub.PutState(shoppieNo, []byte(str))							
	if err != nil {
		return err.Error()
	}
	// validate insert
	//shoppieAsBytes, err = stub.GetState(shoppieNo)
	//if err != nil {
	//	return err.Error()
	//}
	
	return ""
	//res := MyLittleShopper{}
	//json.Unmarshal(shoppieAsBytes, &res)
	// if res.ShoppieNo == shoppieNo {
	// 	return ""
	// }
	// jsonResp = "{\"Error\":\"Failed to get state for " + shoppieNo + " : " + res.ShoppieNo + "\"}"
	// return jsonResp
	//jsonAsBytes, _ := json.Marshal(res)
	//return string(jsonAsBytes)
}

