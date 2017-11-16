// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts  
*/ 
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Tuna structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type Tuna struct {
	Borrower string `json:"borrower"`
	Timestamp string `json:"timestamp"`
	Details  string `json:"details"`
	Lender  string `json:"lender"`
}

/*
 * The Init metd *
 called when the Smart Contract  instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	fmt.Println(function);

	if function == "queryLoan" {
		return s.queryLoan(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordLoan" {
		return s.recordLoan(APIstub, args)
	} else if function == "queryAllLoan" {
		return s.queryAllLoan(APIstub)
	} else if function == "changeLoanLender" {
		return s.changeLoanLender(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryLoan method *
Used to view the records of one particular Loan
It takes one argument -- the key for the Loan in question
 */
func (s *SmartContract) queryLoan(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	loanAsBytes, _ := APIstub.GetState(args[0])
	if loanAsBytes == nil {
		return shim.Error("Could not locate loan")
	}
	return shim.Success(loanAsBytes)
}

/*
 * The initLedger method *
Will add test data (10 loan catches)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	loan := []Tuna{
		Tuna{Borrower: "456789", Details: "67.0006, -70.5476", Timestamp: "1504054225", Lender: "Miriam"},
		Tuna{Borrower: "M83T", Details: "91.2395, -49.4594", Timestamp: "1504057825", Lender: "Dave"},
		Tuna{Borrower: "T012", Details: "58.0148, 59.01391", Timestamp: "1493517025", Lender: "Igor"},
		Tuna{Borrower: "P490", Details: "-45.0945, 0.7949", Timestamp: "1496105425", Lender: "Amalea"},
		Tuna{Borrower: "S439", Details: "-107.6043, 19.5003", Timestamp: "1493512301", Lender: "Rafa"},
		Tuna{Borrower: "J205", Details: "-155.2304, -15.8723", Timestamp: "1494117101", Lender: "Shen"},
		Tuna{Borrower: "S22L", Details: "103.8842, 22.1277", Timestamp: "1496104301", Lender: "Leila"},
		Tuna{Borrower: "EI89", Details: "-132.3207, -34.0983", Timestamp: "1485066691", Lender: "Yuan"},
		Tuna{Borrower: "129R", Details: "153.0054, 12.6429", Timestamp: "1485153091", Lender: "Carlo"},
		Tuna{Borrower: "49W4", Details: "51.9435, 8.2735", Timestamp: "1487745091", Lender: "Fatima"},
	}

	i := 0
	for i < len(loan) {
		fmt.Println("i is ", i)
		loanAsBytes, _ := json.Marshal(loan[i])
		APIstub.PutState(strconv.Itoa(i+1), loanAsBytes)
		fmt.Println("Added", loan[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordLoan method *
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordLoan(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var loan = Tuna{ Borrower: args[1], Details: args[2], Timestamp: args[3], Lender: args[4] }

	loanAsBytes, _ := json.Marshal(loan)
	err := APIstub.PutState(args[0], loanAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record loan: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllLoan method *
allows for assessing all the records added to the ledger(all Loans)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllLoan(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllLoan:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeLoanLender method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, Loan id and new lender name. 
 */
func (s *SmartContract) changeLoanLender(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	loanAsBytes, _ := APIstub.GetState(args[0])
	if loanAsBytes == nil {
		return shim.Error("Could not locate loan")
	}
	loan := Tuna{}

	json.Unmarshal(loanAsBytes, &loan)
	// Normally check that the specified argument is a valid lender of loan
	// we are skipping this check for this example
	loan.Lender = args[1]

	loanAsBytes, _ = json.Marshal(loan)
	err := APIstub.PutState(args[0], loanAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change Loan lender: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}