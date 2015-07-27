package action

import (
	"common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/garyburd/redigo/redis"
)

type AtStorage struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

func NewATStorage() *AtStorage {
	r := &AtStorage{}

	client := &osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:8080",
		UserData:    "",
	}

	setValue("1234", toJSON(client))

	return r
}

func (s *AtStorage) Clone() osin.Storage {
	return s
}

func (s *AtStorage) Close() {
}

// client
// key
// 	- id
// 	- id:secret
// 	- id:redirecturl

// - authorize:code

func getValue(key string) string {
	conn := common.GetRedisPool().Get()
	ret, err := redis.String(conn.Do("get", key))
	defer conn.Close()
	if err != nil {
		fmt.Println("Method - getValue : ", err)
		return ""
	}
	return ret
}

func setValue(key, value string) bool {
	conn := common.GetRedisPool().Get()
	_, err := conn.Do("set", key, value)
	defer conn.Close()
	if err != nil {
		fmt.Println("Method - setValue : ", err)
		return false
	}
	return true
}

func delValue(key string) bool {
	conn := common.GetRedisPool().Get()
	_, err := conn.Do("del", key)
	defer conn.Close()
	if err != nil {
		fmt.Println("Method - delValue : ", err)
		return false
	}
	return true
}

func fromJSON(jsonBytes string, obj interface{}) bool {
	err := json.Unmarshal([]byte(jsonBytes), &obj)
	if err != nil {
		fmt.Println("Methdo - fromJSON : ", err)
		return false
	}
	return true
}

func toJSON(obj interface{}) string {
	ret, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Methdo - toJSON : ", err)
		return ""
	}
	return string(ret)
}

func (s *AtStorage) GetClient(id string) (osin.Client, error) {
	fmt.Printf("GetClient: %s\n", id)
	if id == "" {
		return nil, nil
	}
	value := getValue(id)
	if value == "" {
		return nil, errors.New("Client not found")
	}
	var client osin.DefaultClient
	fromJSON(value, &client)
	return &client, nil
}

func (s *AtStorage) SetClient(id string, client osin.Client) error {
	fmt.Printf("SetClient: %s\n", id)
	ret := toJSON(client)
	setValue(id, ret)
	return nil
}

func (s *AtStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	fmt.Printf("SaveAuthorize: %s\n", data.Code)
	ret := toJSON(data)
	setValue("authorize:"+data.Code, ret)
	return nil
}

func (s *AtStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	fmt.Printf("LoadAuthorize: %s\n", code)

	value := getValue("authorize:" + code)
	//fmt.Println(">>>>>>", value)
	if value == "" {
		return nil, errors.New("Authorize not found")
	}
	// var authorizeData osin.AuthorizeData
	var authorizeData AuthorizeData
	fromJSON(value, &authorizeData)
	// fmt.Println(">>>>>>", authorizeData)
	//fmt.Println(">>>>>>", authorizeData.transfer())
	// return &authorizeData, nil
	return authorizeData.transfer(), nil
}

func (s *AtStorage) RemoveAuthorize(code string) error {
	fmt.Printf("RemoveAuthorize: %s\n", code)
	return nil
	ret := delValue("authorize:" + code)
	if ret == false {
		return errors.New("Del Authorize Faild")
	}
	return nil
}

func (s *AtStorage) SaveAccess(data *osin.AccessData) error {
	fmt.Printf("SaveAccess:%s\n", data.AccessToken)
	ret := toJSON(data)
	strKey := "access:" + data.AccessToken
	setValue(strKey, ret)
	if data.RefreshToken != "" {
		//		setValue("refresh:"+data.RefreshToken, toJSON(data.AccessToken))
		setValue("refresh:"+data.RefreshToken, data.AccessToken)
	}
	return nil
}

func (s *AtStorage) LoadAccess(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadAccess:%s\n", code)
	strKey := "access:" + code
	ret := getValue(strKey)
	fmt.Println("strKey", strKey)
	fmt.Println("LoadAccess", ret)

	if ret == "" {
		return nil, errors.New("access not found")
	}
	fmt.Printf("LoadAccess:\n%s\n", ret)
	var accessData AccessData
	fromJSON(ret, &accessData)
	return accessData.transfer(), nil
}

func (s *AtStorage) RemoveAccess(code string) error {
	fmt.Printf("RemoveAccess: %s\n", code)

	ret := delValue("access:" + code)
	if ret == false {
		return errors.New("Del access Faild")
	}
	return nil
}

func (s *AtStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	fmt.Printf("LoadRefresh: %s\n", code)

	ret := getValue("refresh:" + code)
	if ret == "" {
		return nil, errors.New("Refresh not found")
	}
	return s.LoadAccess(ret)
}

func (s *AtStorage) RemoveRefresh(code string) error {
	fmt.Printf("RemoveRefresh: %s\n", code)
	ret := delValue("refresh:" + code)
	if ret == false {
		return errors.New("Del Refresh faild")
	}
	return nil
}
