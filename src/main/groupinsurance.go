package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

type GroupPolicy struct {
	ObjectType string `json:"docType"`
	PolicyNo string `json:"policyNo"`
	CustomerId string `json:"customerId"`
	//Insured Insured `json:"insured"`
	Coverages []Coverage `json:"coverages"`
	TransactionType string `json:"transactionType"`
	TransactionLabel string `json:"transactionLabel"`
	TransactionDetails string `json:"transactionDetails"`
}	

type Insured struct {
	ObjectType string `json:"docType"`
	EmployeeId string `json:"employeeId"`
	CustomerId string `json:"customerId"`
	CertificateNo string `json:"certificateNo"`
	Class string `json:"class"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	EmployerId string `json:"employerId"`
	EmployerName string `json:"employerName"`
	Depdendents []Dependent `json:"depdendents"`
}	

type Coverage struct {
	ObjectType string `json:"docType"`
	coverageType string `json:"coverageType"`
	coverageLabel string `json:"coverageLabel"`
	sumAssured string `json:"sumAssured"`
}	

type Dependent struct {
	ObjectType string `json:"docType"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Relatioship string `json:"relatioship"`
}	

/*


*/
func (t *GroupPolicy) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	fmt.Println("Init function start")
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("Init function end")
	
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *GroupPolicy) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "enroll" {
		return t.enroll(stub, args)
	} else if function == "updateClass" {
		return t.updateClass(stub, args)
	} else if function == "addUpdateDependent" {
		return t.addUpdateDependent(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *GroupPolicy) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "findPolicyByPolicyNo" {											//read a variable
		return t.findPolicyByPolicyNo(stub, args)
	} else if function =="findPolicyByEmployer" {
		return t.findPolicyByEmployer(stub, args)
	} else if function =="findPolicyByDateRange" {
		return t.findPolicyByDateRange(stub, args)
	} 
	
	fmt.Println("query did not find func: " + function)						//error
	return nil, errors.New("Received unknown function query: " + function)
}

func (t *GroupPolicy) enroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var err error

	fmt.Println("runnin write()")
	
	gp:= new(GroupPolicy)
	gp.ObjectType = "GP"
	gp.PolicyNo = args[0]
	gp.CustomerId = args[1]
	gp.TransactionType = args[2]
	gp.TransactionLabel = args[3]
	gp.TransactionDetails = args[4]
//	
//	insured:= new(Insured)
//	insured.ObjectType="INS"
//	insured.CustomerId = args[5]
//	insured.EmployeeId = args[6]
//	insured.FirstName = args[7]
//	insured.LastName = args[8]
//	insured.CertificateNo = args[9]
//	insured.Class = args[10]
//	insured.EmployerId = args[11]
//	insured.EmployerName = args[12]
	
	//gp.Insured = insured

	jsonAsBytes, _ := json.Marshal(gp) 
	
	//fmt.Println("jsonAsBytes >> " + jsonAsBytes)
	
	err = stub.PutState(args[0], jsonAsBytes)
	 
	if err != nil {
        return nil, err
    }
    
    return nil, nil
}

func (t *GroupPolicy) updateClass(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	
	fmt.Println("runnin write()")
	key = args[0]
	value = args[1]

	err = stub.PutState(key, []byte(value))
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number od arguments.");
	}
	
	if err != nil {
        return nil, err
    }
    return nil, nil
}

func (t *GroupPolicy) addUpdateDependent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	
	fmt.Println("runnin write()")
	key = args[0]
	value = args[1]

	err = stub.PutState(key, []byte(value))
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number od arguments.");
	}
	
	if err != nil {
        return nil, err
    }
    return nil, nil
}

func (t *GroupPolicy) findPolicyByPolicyNo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var policyNo, jsonResp string
	var err error
	if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }
	policyNo = args[0]
	valAsbytes, err := stub.GetState(policyNo)
	
	if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + policyNo + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil

}

func (t *GroupPolicy) findPolicyByEmployeeId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }
	valAsbytes, err := stub.GetState(key)
	
	if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil

}

func (t *GroupPolicy) findPolicyByEmployer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }
	valAsbytes, err := stub.GetState(key)
	
	if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil

}

func (t *GroupPolicy) findPolicyByDateRange(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }
	valAsbytes, err := stub.GetState(key)
	
	if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil

}

func main() {
	err := shim.Start(new(GroupPolicy))
	if err != nil {
		fmt.Printf("Error starting Claim: %s", err)
	}
} 

