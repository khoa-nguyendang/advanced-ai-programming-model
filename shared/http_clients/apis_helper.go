package httpclients

import (
	"aapi/shared/constants"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func SendHttpGet(url string,
	requestType constants.HttpRequestEnum,
	responseModel interface{},
	token, username, password string,
) error {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	if requestType == constants.Bear {
		req.Header.Set("Authorization", "bearer "+token)
	}

	if requestType == constants.Basic {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return json.Unmarshal(body, &responseModel)
}

func SendHttpPost(url string,
	requestType constants.HttpRequestEnum,
	payload interface{},
	responseModel interface{},
	token, username, password string,
) error {

	payloadData, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadData))
	req.Header.Set("Content-Type", "application/json")

	if requestType == constants.Bear {
		req.Header.Set("Authorization", "bearer "+token)
	}

	if requestType == constants.Basic {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	log.Printf("body: %v", string(body))
	switch v := responseModel.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		log.Printf("responseModel was Integer: %v", v)
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		log.Printf("responseModel was  Float64: %v", v)
	case string:
		log.Printf("responseModel was string: %v", v)
		responseModel = string(body)
	default:
		log.Printf("responseModel was struct: %v", v)
		return json.Unmarshal(body, &responseModel)
	}

	return nil
}

func SendHttpPost2(url string,
	requestType constants.HttpRequestEnum,
	payload interface{},
	token, username, password string,
) ([]byte, error) {

	payloadData, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	log.Println("payloadData: %v", string(payloadData))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadData))
	req.Header.Set("Content-Type", "application/json")

	if requestType == constants.Bear {
		req.Header.Set("Authorization", "bearer "+token)
	}

	if requestType == constants.Basic {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
