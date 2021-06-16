package server

import (
	"GoMS1Assignment/restapi/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kennygrant/sanitize"
)

// courseInfo struct for the json
type courseInfo struct {
	Title string `json:"Title"`
}

// validKey checks that the access keys supplied to the REST API is valid. Returns a bool.
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		urlKey := sanitize.Accents(key[0]) // Sanitise the url param string

		// Checks key in the database
		if database.ValidKey(urlKey) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// home is the handler function to display the home page of the REST API.
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the GoSchool REST API!")
}

// allcourses is the handler function to retrieve all courses.
// It converts the map object retrieved into JSON and passes it back to the client.
func allcourses(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	// Get all courses from the database
	courses := database.GetAllCourses()

	// convert the map object to JSON, and pass it back to the client
	json.NewEncoder(w).Encode(courses)
}

// courses is the handler function for CRUD operations sent by the client.
// The operations for GET, POST, PUT and DELETE will be determined by switch
// and its respective functions will be called.
func course(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)

	switch r.Method {
	case "GET": // GET is for retrieving course
		getCourse(params, w)
	case "POST": // POST is for creating new course
		addCourse(params, w, r)
	case "PUT": //---PUT is for updating course
		updateCourse(params, w, r)
	case "DELETE": // DELETE is for deleting course
		deleteCourse(params, w)
	}
}

// getCourse implements the GET method invoked by the client and
// retrieves the course detail with the course id given.
func getCourse(params map[string]string, w http.ResponseWriter) {
	// Get courses from the database
	courses := database.GetCourse(params["courseid"])

	// Course exists
	if len(courses) != 0 {
		// convert the map object to JSON, and pass it back to the client
		json.NewEncoder(w).Encode(courses)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No course found"))
	}
}

// addCourse implements the POST method invoked by the client and
// adds a course with the course id and title given.
func addCourse(params map[string]string, w http.ResponseWriter, r *http.Request) {
	// read the string sent to the service
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		newCourse := convertJSON(reqBody, w, r)

		// Check if course exists;
		courses := database.GetCourse(params["courseid"])

		// Course does not exists
		if len(courses) == 0 {
			// Add course information into the database
			database.AddCourse(params["courseid"], newCourse.Title)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Course added: " + params["courseid"]))
		} else {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("409 - Duplicate course ID"))
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply course information in JSON format"))
	}
}

// updateCourse implements the PUT method invoked by the client and
// update a course title with the course id given.
func updateCourse(params map[string]string, w http.ResponseWriter, r *http.Request) {
	// read the string sent to the service
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		newCourse := convertJSON(reqBody, w, r)

		// Check if course exists;
		courses := database.GetCourse(params["courseid"])

		// Course exists
		if len(courses) != 0 {
			// Update course in the database
			database.UpdateCourse(params["courseid"], newCourse.Title)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 - Course updated"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply course information in JSON format"))
	}
}

// deleteCourse implements the DELETE method invoked by the client and
// deletes a course with the course id given.
func deleteCourse(params map[string]string, w http.ResponseWriter) {
	// Check if course exists;
	courses := database.GetCourse(params["courseid"])

	// Course exists
	if len(courses) != 0 {
		// Delete course from the database
		database.DeleteCourse(params["courseid"])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - Course deleted"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No course found"))
	}
}

// convertJSON converts the client JSON to object and returns the unmarshal data.
func convertJSON(reqBody []byte, w http.ResponseWriter, r *http.Request) courseInfo {
	var newCourse courseInfo

	// convert JSON to object
	json.Unmarshal(reqBody, &newCourse)

	if newCourse.Title == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply course information in JSON format"))
	}
	return newCourse
}
