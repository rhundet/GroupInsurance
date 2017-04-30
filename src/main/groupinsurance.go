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
	Insured Insured `json:"insured"`
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
	
	if args[0] == "" {
		return nil, errors.New("PolicyNo missing as argument 0.")
	}
	gp.PolicyNo = args[0]
	
	if args[1] == "" {
		return nil, errors.New("CustomerId missing as argument 1.")
	}
	gp.CustomerId = args[1]
	
	if args[2] == "" {
		return nil, errors.New("TransactionType missing as argument 2.")
	}
	gp.TransactionType = args[2]
	
	if args[3] == "" {
		return nil, errors.New("TransactionLabel missing as argument 3.")
	}
	gp.TransactionLabel = args[3]
	
	if args[4] == "" {
		return nil, errors.New("TransactionDetails missing as argument 4.")
	}
	gp.TransactionDetails = args[4]

	gp.Insured.ObjectType="INS"
	
	if args[5] == "" {
		return nil, errors.New("CustomerId missing as argument 5.")
	}
	gp.Insured.CustomerId = args[5]
	
	if args[6] == "" {
		return nil, errors.New("EmployeeId missing as argument 6.")
	}
	gp.Insured.EmployeeId = args[6]
	
	if args[7] == "" {
		return nil, errors.New("FirstName missing as argument 7.")
	}
	gp.Insured.FirstName = args[7]
	
	if args[8] == "" {
		return nil, errors.New("LastName missing as argument 8.")
	}
	gp.Insured.LastName = args[8]
	
	if args[9] == "" {
		return nil, errors.New("CertificateNo missing as argument 9.")
	}
	gp.Insured.CertificateNo = args[9]
	
	if args[10] == "" {
		return nil, errors.New("Class missing as argument 10.")
	}
	gp.Insured.Class = args[10]
	
	if args[11] == "" {
		return nil, errors.New("EmployerId missing as argument 11.")
	}
	gp.Insured.EmployerId = args[11]
	
	if args[12] == "" {
		return nil, errors.New("EmployerName missing as argument 12.")
	}
	gp.Insured.EmployerName = args[12]
	
	jsonAsBytes, _ := json.Marshal(gp) 
	
	err = stub.PutState(gp.PolicyNo, jsonAsBytes) 
	 
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

