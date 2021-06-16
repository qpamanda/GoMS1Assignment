package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "http://localhost:5000/api/v1/courses"
const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

// addCourse invokes the POST method to the REST API to
// add a course with the course id and JSON data for title given.
// Returns the response code received from the REST API.
func addCourse(code string, jsonData map[string]string) int {
	responseCode := 0
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(baseURL+"/"+code+"?key="+key,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		responseCode = response.StatusCode
		response.Body.Close()
	}
	return responseCode
}

// updateCourse invokes the PUT method to the REST API to
// update a course with the course id and JSON data for title given.
// Returns the response code received from the REST API.
func updateCourse(code string, jsonData map[string]string) int {
	responseCode := 0
	jsonValue, _ := json.Marshal(jsonData)

	request, _ := http.NewRequest(http.MethodPut,
		baseURL+"/"+code+"?key="+key,
		bytes.NewBuffer(jsonValue))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		responseCode = response.StatusCode
		response.Body.Close()
	}
	return responseCode
}

// deleteCourse invokes the DELETE method to the REST API to
// delete a course with the course id given.
// Returns the response code received from the REST API.
func deleteCourse(code string) int {
	responseCode := 0

	request, _ := http.NewRequest(http.MethodDelete, baseURL+"/"+code+"?key="+key, nil)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		responseCode = response.StatusCode
		response.Body.Close()
	}
	return responseCode
}

// getCourse invokes the GET method to the REST API to
// retrieve a course with the course id given. If course id
// is blank, all courses will be retrieved.
// Returns the map object after unmarshalling the JSON received from the REST API.
func getCourse(code string) map[string]courseInfo {
	var courses map[string]courseInfo

	url := baseURL
	if code != "" {
		url = baseURL + "/" + code
	}
	url = url + "?key=" + key

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == 200 { // 200 - http.StatusOk
			json.Unmarshal(data, &courses)
		}
		response.Body.Close()
	}
	return courses
}
