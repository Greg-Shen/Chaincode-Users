package smartcontract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing archives
type SmartContract struct {
	contractapi.Contract
}

// User Data struct
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	fmt.Println("function UserExists")

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJSON != nil, nil
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, id string, name string, email string) error {
	fmt.Println("function CreateUser")
	//檢查是否已經有key
	exists, err := s.UserExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", id)
	}
	//初始化新的user structure
	user := User{
		ID:    id,
		Name:  name,
		Email: email,
	}
	//確認是否屬於json格式
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	ctx.GetStub().PutState(id, userJson) //以user id為key, 把userJson資料傳入

	return nil
}

func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	fmt.Println("function GetUser")

	userJson, err := ctx.GetStub().GetState(id) //將key為id的資料儲存進userJson
	//若錯誤不為空，回傳讀取失敗
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	//若取得資料為空，回傳欲取得資料不存在
	if userJson == nil {
		return nil, fmt.Errorf("the user %s does not exist", id)
	}
	//解析取得資料
	var user User

	err = json.Unmarshal(userJson, &user) //將鏈上的資料格式反序列化
	if err != nil {
		return nil, err
	}

	return &user, nil //若無問題，將資料回傳
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, name string, email string) error {
	fmt.Println("function UpdateUser")
	//此function不可變動id
	user, err := s.GetUser(ctx, id)
	//判斷資料是否存在
	if err != nil {
		return err
	}
	//轉換成輸入欲轉換成的email及名字
	user.Email = email
	user.Name = name

	userJson, err := json.Marshal(user) //將資料序列化
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, userJson) //透過PutState將id,資料存放到DB中
}

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	fmt.Println("function DeleteUser")

	exists, err := s.UserExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the user %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	fmt.Println("function GetAllUsers")

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "") //抓取一個範圍的資料，參數為key的範圍，回傳後為iterator的型態
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close() //close掉iterator

	var users []*User //宣告空array
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next() //把下一個iterator的response抓回來
		if err != nil {
			return nil, err
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user) //透過json.Unmarshal轉成user的structure
		if err != nil {
			return nil, err
		}
		users = append(users, &user) //append到user的陣列
	} //當for迴圈跑完，所有資料會被append完

	return users, nil
}
