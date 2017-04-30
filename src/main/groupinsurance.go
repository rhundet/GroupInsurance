package main


// https://github.com/hyperledger/fabric/blob/v0.6/core/chaincode/shim/chaincode.go
// https://hyperledger-fabric.readthedocs.io/en/v0.6/API/ChaincodeAPI.html

//http://hyperledger-fabric.readthedocs.io/en/stable/API/ChaincodeAPI/#chaincode-apis
// https://github.com/hyperledger/fabric/blob/master/core/chaincode/shim/chaincode.go

/**
Not implementing tables as I couldn't find the tables API in 1.0 or latest API so want to keep code portable to 1.0
Can't implement queries as not present 1.0 API and I am using 0.6 on bluemix.
Therefore keeping the key as employee or customer Id which will be same for the network
and be accessible by all peers

**/
import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	//"strconv"
)

const MAX_NO_OF_DEPENDENTS = 3

type GroupPolicy struct {
	ObjectType string `json:"docType"`
	PolicyNo string `json:"policyNo"`
	CustomerId string `json:"customerId"` // this will be our customerId as well as employeeId
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
	if function == "findPolicyByEmployeeId" {			
		employeeId:= args[0]	//read a variable
		return t.findPolicyByEmployeeId(stub, employeeId)
	} 
	
	fmt.Println("query did not find func: " + function)						//error
	return nil, errors.New("Received unknown function query: " + function)
}

/**
This will create a new customer or enroll a new customer
*/
func (t *GroupPolicy) enroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var err error

	fmt.Println("runnin write()")
	
	if len(args) != 13 {
		return nil, errors.New("Incorrect number of arguments. Expecting 13")
	}

	gp:= new(GroupPolicy)
	gp.ObjectType = "GP"
	gp.PolicyNo = args[0]
	gp.CustomerId = args[1] // this should be employeeId as well
	gp.TransactionType = args[2]
	gp.TransactionLabel = args[3]
	gp.TransactionDetails = args[4]

	gp.Insured.ObjectType="INS"
	gp.Insured.CustomerId = args[5]
	gp.Insured.EmployeeId = args[6] // This should be customerId
	gp.Insured.FirstName = args[7]
	gp.Insured.LastName = args[8]
	gp.Insured.CertificateNo = args[9]
	gp.Insured.Class = args[10]
	gp.Insured.EmployerId = args[11]
	gp.Insured.EmployerName = args[12]
	
	jsonAsBytes, _ := json.Marshal(gp) 
	
	err = stub.PutState(gp.Insured.EmployeeId, jsonAsBytes)
	 
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
	
	var employeeId string
	var err error
	
	fmt.Println("runnin write()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number od arguments.");
	}

	employeeId = args[0]
	currentPolicyJson,_:= t.findPolicyByEmployeeId(stub, employeeId)
	currentPolicy:= new(GroupPolicy);
	json.Unmarshal(currentPolicyJson, &currentPolicy) 
	
	// current number of dependents
	noOfDependents:= len(currentPolicy.Insured.Depdendents)
	
	// maximum 3 dependents can be added
	if noOfDependents < 3 {
		// add new dependents
		newCount:= noOfDependents + 1
		
		//back up current dependents
		oldDependents:= currentPolicy.Insured.Depdendents
		currentPolicy.Insured.Depdendents = make([]Dependent,(newCount))
		
		// copy old dependents as it is
		for i := 0; i < noOfDependents; i++ {
			currentPolicy.Insured.Depdendents[i].ObjectType = oldDependents[i].ObjectType
			currentPolicy.Insured.Depdendents[i].FirstName = oldDependents[i].FirstName
			currentPolicy.Insured.Depdendents[i].LastName = oldDependents[i].LastName
			currentPolicy.Insured.Depdendents[i].Relatioship = oldDependents[i].Relatioship
		}
		
		// create new dependent
		lastIndex:= noOfDependents;
		currentPolicy.Insured.Depdendents[lastIndex].ObjectType = "DEP"
		currentPolicy.Insured.Depdendents[lastIndex].FirstName = args[1]
		currentPolicy.Insured.Depdendents[lastIndex].LastName = args[2]
		currentPolicy.Insured.Depdendents[lastIndex].Relatioship = args[3]
	}
	
	updatedJsonAsBytes, _ := json.Marshal(currentPolicy)
	err = stub.PutState(employeeId, updatedJsonAsBytes)
	
	if err != nil {
        return nil, err
    }
    return nil, nil
}

func (t *GroupPolicy) findPolicyByEmployeeId(stub shim.ChaincodeStubInterface, employeeId string) ([]byte, error) {
	var jsonResp string
	var err error
	
	valAsbytes, err := stub.GetState(employeeId)
	
	if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + employeeId + "\"}"
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

