package main

// 외부모듈
import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SmartContract 클래스 정의
type SmartContract struct {
}

// Member. 구조체 정의
type Member struct {
	MemberNo            string `json:"memberno"`
	MemberId            string `json:"memberid"`
	MemberName          string `json:"membername"`
	MemberPWD           string `json:"memberpwd"`
	MemberToEthereumKey string `json:"membertoethereumkey"`
	MemberEntryDate     string `json:"memberentrydate"`
}

// Init 함수 구현
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// 초기 시스템의 상태 설정
	// instantiate 권한 ca에 등록된 배포한사람의 role관리
	return shim.Success(nil)
}

// Invoke 함수 구현
// peer chaincode invoke -n Member -C mychannel -c '{"Args":[]}'
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()

	if fn == "memberc" { // 블록생성
		return s.memberc(stub, args)
	} else if fn == "memberu" { // 블록수정
		return s.memberu(stub, args)
	} else if fn == "memberd" { // 블록제거
		return s.memberd(stub, args)
	} else if fn == "query" { // ws 조회
		return s.query(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// register //  Creation R U D
// peer chaincode invoke -n Member -C mychannel -c '{"Args":[]}
func (s *SmartContract) memberc(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of argument. Expecting 6")
	}

	// (TO DO) 이미 등록된 EVENT ID인가?
	memberAsBytes, _ := stub.GetState(args[0])
	if memberAsBytes != nil {
		return shim.Error("Event id already exists.")
	}

	var member = Member{MemberNo: args[0], MemberId: args[1], MemberName: args[2], MemberPWD: args[3], MemberToEthereumKey: args[4], MemberEntryDate: args[5]}

	memberAsBytes, _ = json.Marshal(member)
	stub.PutState(args[0], memberAsBytes)

	return shim.Success(nil)
}

// join //  C R Update D
// peer chaincode invoke -n Member -C mychannel -c '{"Args":[]}' ('이름, 패스워드, 이더리움주소'만 수정)
func (s *SmartContract) memberu(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of argument. Expecting 5")
	}

	// World State 조회
	memberAsBytes, _ := stub.GetState(args[0])
	if memberAsBytes == nil {
		return shim.Error("Requested member id is missing")
	}

	var member = Member{MemberId: args[0], MemberName: args[1], MemberPWD: args[2], MemberToEthereumKey: args[3]}

	memberAsBytes, _ = json.Marshal(member)
	stub.PutState(args[0], memberAsBytes)

	return shim.Success(nil)
}

// finalize //  C R Update D
// peer chaincode invoke -n member -C mychannel -c '{"Args":[]}'
func (s *SmartContract) memberd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of argument. Expecting 1")
	}

	memberAsBytes, _ := stub.GetState(args[0])
	if memberAsBytes == nil {
		return shim.Error("Requested member id is missing")
	}

	stub.DelState(args[0])

	return shim.Success(nil)
}

// query //  C Rear/Retrieve U D
// peer chaincode query -n member -C mychannel -c '{"Args":[]}'
func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of argument. Expecting 1")
	}

	// World State 조회
	memberAsBytes, _ := stub.GetState(args[0])
	if memberAsBytes == nil {
		return shim.Error("Requested member id is missing")
	}
	return shim.Success(memberAsBytes)
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// main
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
