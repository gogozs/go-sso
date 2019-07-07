package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


type Response struct {
	Code int `json:"code"`
	Data Data `json:"Data"`
}

type Data struct {
	Total int `json:"total"`
	Items []Items `json:"items"`
}

type Items struct {
	AssetId string `json:"asset_id"`
	AssetType string `json:"asset_type"`
}


func GetAsset() ([]Items, error) {
	var data Response
	req, err := http.NewRequest("GET", "", nil)

	client := &http.Client{}
	req.Header.Add("Authorization", "Token e7f4519ec7cce3536258165dc262488eba3e7ad0")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body, &data)
	fmt.Println(data)
	return data.Data.Items, err
}


