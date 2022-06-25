package creatiointegrationgolang

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Result struct {
	success        bool
	message        string
	responseStatus int
	response       []byte
}

var BPMCSRF = &http.Cookie{}
var ASPXAUTH = &http.Cookie{}
var BPMLOADER = &http.Cookie{}
var UserNameC = &http.Cookie{}
var URL string
var UserName string
var UserPassword string

/**
* Method rest
* Creatio RESTFull Request
* @param {map[string]string} endpoint
*              endpoint.service {string}
*              endpoint.method {string}
* @param {[]byte](Json)} data
* @return {object} result
*              	result.success {bool}
*              	result.message {string}
*				result.responseStatus {int}
*              	result.response {[]byte](Json)}
 */
func Rest(method string, endpoint map[string]string, data []byte) Result {
	var result Result
	initRequest := getAuth(UserName, UserPassword)
	if initRequest.success {
		requestRest := getRestRequest(method, endpoint, data)
		result = requestRest
	} else {
		result = initRequest
	}

	return result
}

/**
*
 */
func getAuth(userName string, userPassword string) Result {

	var result Result

	var credential = map[string]string{"UserName": userName, "UserPassword": userPassword}
	bodyRequest, _ := json.Marshal(credential)

	response, err := http.Post(URL+"/ServiceModel/AuthService.svc/Login", "application/json", bytes.NewBuffer(bodyRequest))

	if err != nil {
		result.message = err.Error()
		result.success = false
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		result.message = err.Error()
		result.success = false
	}

	defer response.Body.Close()

	cookies := response.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "BPMCSRF" {
			BPMCSRF = cookie
		}
		if cookie.Name == ".ASPXAUTH" {
			ASPXAUTH = cookie
		}
		if cookie.Name == "BPMLOADER" {
			BPMLOADER = cookie
		}
		if cookie.Name == "UserName" {
			UserNameC = cookie
		}

	}

	result.response = responseBody
	result.responseStatus = response.StatusCode
	result.success = true

	return result

}

/**
*
 */
func getRestRequest(method string, endpoint map[string]string, data []byte) Result {

	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
		Jar:     jar,
	}
	var result Result
	url := URL + "/0/rest/" + endpoint["service"] + "/" + endpoint["method"]

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		result.message = err.Error()
		result.success = false
	}

	req.Header.Set("BPMCSRF", string(BPMCSRF.Value))
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(ASPXAUTH)
	req.AddCookie(BPMCSRF)
	req.AddCookie(BPMLOADER)
	req.AddCookie(UserNameC)

	response, err := client.Do(req)

	if err != nil {
		result.message = err.Error()
		result.success = false
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		result.message = err.Error()
		result.success = false
	}

	defer response.Body.Close()

	result.response = responseBody
	result.responseStatus = response.StatusCode
	result.success = true

	return result
}
