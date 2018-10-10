package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Display struct {
	ScheduleURL  string `json:"scheduleURL"`
	ScheduleHash  string `json:"scheduleHash"`
}

/*
 * The Init method is called when the Smart Contract "scheduler" is instantiated by the blockchain network
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "scheduler"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryDisplay" {
		return s.queryDisplay(APIstub, args)
	} else if function == "createOrUpdateDisplay" {
		return s.createOrUpdateDisplay(APIstub, args)
	} else if function == "updateScheduleHash" {
		return s.updateScheduleHash(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryDisplay(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	displayAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(displayAsBytes)
}

func (s *SmartContract) createOrUpdateDisplay(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var display = Display{ScheduleURL: args[1], ScheduleHash: args[2]}

	displayAsBytes, _ := json.Marshal(display)
	APIstub.PutState(args[0], displayAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) updateScheduleHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	displayAsBytes, _ := APIstub.GetState(args[0])
	display := Display{}

	json.Unmarshal(displayAsBytes, &display)
        display.ScheduleHash = args[1]

	displayAsBytes, _ = json.Marshal(display)
	APIstub.PutState(args[0], displayAsBytes)

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
