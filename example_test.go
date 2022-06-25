package creatiointegrationgolang

import (
	"fmt"
	"testing"
)

func TestExampleRest(t *testing.T) {
	UserName = "Superadmin"
	UserPassword = "Superadmin"
	URL = "http://localhost:8080"

	var JsonString = `{ "Request" :{"Month" :"March","Year" :"2022","Company" :"","Directorate":"","Department" :""}}`
	var JsonData = []byte(JsonString)

	endpoint := map[string]string{
		"service": "CustomDashboardAPI",
		"method":  "GetSumberDana",
	}

	response := Rest("POST", endpoint, JsonData)

	fmt.Println(response)
	//do somthing Here

}
