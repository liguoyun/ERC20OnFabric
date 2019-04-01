package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

const TokenName = "ERC20Token"
const TokenSym = "ERC20"
const CoinBase = "CoinBase"

type ERC20TokenChaincode struct {
}

func (s *ERC20TokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	//debug purpose,to check if the new chaincode updated
	fmt.Printf("Init Token V1.3 \n")
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 for Init")
	}
	owner := args[0]
	totalSuplly, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error("Total Supply conversion error, interger expected")
	}
	token := NewERC20Token(TokenName, TokenSym, owner, totalSuplly)

	tokenAsBytes, err := json.Marshal(token)
	err = stub.PutState(TokenName, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Init Token %s \n", string(tokenAsBytes))
	}
	return shim.Success(nil)
}
func (s *ERC20TokenChaincode) totalSupply(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tokenAsBytes, err := stub.GetState(TokenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := ERC20Token{}
	json.Unmarshal(tokenAsBytes, &token)
	totalStr := strconv.FormatInt(token.TotalSupply, 10)
	fmt.Println("TotalSupply is ", token.TotalSupply)
	return shim.Success([]byte(totalStr))
}

func (s *ERC20TokenChaincode) name(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tokenAsBytes, err := stub.GetState(TokenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := ERC20Token{}
	json.Unmarshal(tokenAsBytes, &token)
	fmt.Println("Token Name is ", token.TokenName)
	return shim.Success([]byte(token.TokenName))
}

func (s *ERC20TokenChaincode) symbol(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tokenAsBytes, err := stub.GetState(TokenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := ERC20Token{}
	json.Unmarshal(tokenAsBytes, &token)
	fmt.Println("Token Symbol is ", token.TokenSymbol)
	return shim.Success([]byte(token.TokenSymbol))
}

//transfer from owner to a
func (s *ERC20TokenChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	to := args[0]
	amount, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error(fmt.Errorf("Convert amount to int64 failed").Error())
	}
	fmt.Println("Transfer  ", amount, " to  ", to)
	tokenAsBytes, err := stub.GetState(TokenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := ERC20Token{}
	json.Unmarshal(tokenAsBytes, &token)

	err = token.transfer(token.TokenOwner, to, int64(amount))
	if err != nil {
		return shim.Error(err.Error())
	}
	tokenBytes, err := json.Marshal(token)
	err = stub.PutState(TokenName, tokenBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Save Token %s \n", string(tokenBytes))
	}
	return shim.Success(nil)
}

//transfer from a to b
func (s *ERC20TokenChaincode) transferFrom(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	from := args[0]
	to := args[1]
	amount, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error(fmt.Errorf("Convert amount to int64 failed").Error())
	}
	fmt.Println("Transfer  ", amount, " from ", from, " to  ", to)
	tokenAsBytes, err := stub.GetState(TokenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := ERC20Token{}
	json.Unmarshal(tokenAsBytes, &token)

	err = token.transfer(from, to, int64(amount))
	if err != nil {
		return shim.Error(err.Error())
	}
	tokenBytes, err := json.Marshal(token)
	err = stub.PutState(TokenName, tokenBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Save Token %s \n", string(tokenBytes))
	}
	return shim.Success(nil)
}

func (s *ERC20TokenChaincode) banalanceOf(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	tokenAsBytes, err := stub.GetState(TokenName)
	if err != nil {
		return shim.Error(err.Error())
	}
	token := ERC20Token{}
	json.Unmarshal(tokenAsBytes, &token)
	amount := token.Balances[args[0]]
	amountStr := strconv.FormatInt(amount, 10)
	fmt.Println(args[0], "Balance is ", amountStr)
	return shim.Success([]byte(amountStr))
}

func (s *ERC20TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// Retrieve the requested Chaincode function and arguments
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(function, " was called")
	if function == "name" {
		return s.name(stub, args)
	} else if function == "symbol" {
		return s.symbol(stub, args)
	} else if function == "totalSupply" {
		return s.totalSupply(stub, args)
	} else if function == "balanceOf" {
		return s.banalanceOf(stub, args)
	} else if function == "transfer" {
		return s.transfer(stub, args)
	} else if function == "transferFrom" {
		return s.transferFrom(stub, args)
	}
	return shim.Error("Invalid Chaincode function name.")
}

func main() {
	err := shim.Start(new(ERC20TokenChaincode))
	if err != nil {
		fmt.Printf("Error starting ERC20Token chaincode: %s", err)
	}
}
