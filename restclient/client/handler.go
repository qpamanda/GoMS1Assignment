package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kennygrant/sanitize"
)

// courseInfo struct for the json
type courseInfo struct {
	Title string `json:"Title"`
}

// index is the handler function to display the home page of the client.
// This is also the page where all courses are retrieved and displayed.
// CRUD function for getCourse is invoked.
func index(w http.ResponseWriter, r *http.Request) {
	courses := getCourse("") // Get all courses
	tpl.ExecuteTemplate(w, "index.gohtml", courses)
}

// addcourse is the handler function to retrieve user input for new course details.
// Validations are performed to ensure valid course details are submitted.
// CRUD function for addCourse is invoked.
func addcourse(w http.ResponseWriter, r *http.Request) {
	clientMsg := "" // To display message to the user on the client
	courseID := ""
	courseTitle := ""

	if r.Method == http.MethodPost {
		courseID = r.FormValue("courseid")
		courseTitle = r.FormValue("coursetitle")

		if err := validateCourseID(courseID); err != nil {
			clientMsg = err.Error()
		} else if err := validateCourseTitle(courseTitle); err != nil {
			clientMsg = err.Error()
		} else {
			jsonData := map[string]string{"title": courseTitle}

			responseCode := addCourse(courseID, jsonData)

			if responseCode == 201 { // 201 - http.StatusCreated
				clientMsg = fmt.Sprintf("%s - %s added successfully.\n", courseID, courseTitle)
			} else if responseCode == 409 { // 409 - http.StatusConflict
				clientMsg = ">> Duplicate Course ID."
			} else {
				clientMsg = ">> Error adding course. Please contact the system administrator."
			}
		}
	}

	data := struct {
		CourseID    string
		CourseTitle string
		ClientMsg   string
	}{
		courseID,
		courseTitle,
		clientMsg,
	}

	tpl.ExecuteTemplate(w, "addcourse.gohtml", data)
}

// updcourse is the handler function to retrieve user input for change in course title.
// Validations are performed to ensure valid course title are submitted.
// CRUD function for updateCourse is invoked.
func updcourse(w http.ResponseWriter, r *http.Request) {
	clientMsg := "" // To display message to the user on the client
	courseID := ""
	courseTitle := ""
	validCourseID := true // Determine whether to show course info

	v := r.URL.Query()
	if key, ok := v["courseid"]; ok {
		courseID = sanitize.Accents(key[0])
	}

	if courseID != "" {
		courses := getCourse(courseID) // Get all courses
		if len(courses) == 0 {
			validCourseID = false
			clientMsg = ">> Invalid Course ID"
		} else {
			courseTitle = courses[courseID].Title
		}
	}

	if r.Method == http.MethodPost {
		courseID = r.FormValue("courseid")
		courseTitle = r.FormValue("coursetitle")

		if err := validateCourseTitle(courseTitle); err != nil {
			clientMsg = err.Error()
		} else {
			jsonData := map[string]string{"title": courseTitle}
			responseCode := updateCourse(courseID, jsonData)

			if responseCode == 200 { // 200 - http.StatusOK
				clientMsg = fmt.Sprintf("%s - %s updated successfully.\n", courseID, courseTitle)
			} else if responseCode == 404 { // 404 - http.StatusNotFound
				clientMsg = ">> Course not found."
			} else {
				clientMsg = ">> Error updating course."
			}
		}
	}

	data := struct {
		CourseID      string
		CourseTitle   string
		ClientMsg     string
		ValidCourseID bool
	}{
		courseID,
		courseTitle,
		clientMsg,
		validCourseID,
	}

	tpl.ExecuteTemplate(w, "updcourse.gohtml", data)
}

// delcourse is the handler function to delete a course details as selected by the user.
// CRUD function for deleteCourse is invoked.
func delcourse(w http.ResponseWriter, r *http.Request) {
	clientMsg := "" // To display message to the user on the client
	courseID := ""
	courseTitle := ""
	validCourseID := true // Determine whether to show course info

	v := r.URL.Query()
	if key, ok := v["courseid"]; ok {
		courseID = sanitize.Accents(key[0])
	}

	if courseID != "" {
		courses := getCourse(courseID) // Get all courses
		if len(courses) == 0 {
			validCourseID = false
			clientMsg = ">> Invalid Course ID"
		} else {
			courseTitle = courses[courseID].Title
		}
	}

	if r.Method == http.MethodPost {
		courseID = r.FormValue("courseid")
		courseTitle = r.FormValue("coursetitle")

		responseCode := deleteCourse(courseID)

		if responseCode == 200 { // 200 - http.StatusOK
			clientMsg = fmt.Sprintf("%s - %s deleted successfully.\n", courseID, courseTitle)
		} else if responseCode == 404 { // 404 - http.StatusNotFound
			clientMsg = ">> Course not found."
		} else {
			clientMsg = ">> Error deleting course."
		}
	}

	data := struct {
		CourseID      string
		CourseTitle   string
		ClientMsg     string
		ValidCourseID bool
	}{
		courseID,
		courseTitle,
		clientMsg,
		validCourseID,
	}

	tpl.ExecuteTemplate(w, "delcourse.gohtml", data)
}

// validateCourseID checks that user input for course id is valid. Returns error type.
func validateCourseID(courseID string) error {
	if courseID != "" {
		if len(courseID) > 6 {
			return errors.New(">> Course ID cannot be greater than 6 characters")
		}
	} else {
		return errors.New(">> Course ID cannot be blank")
	}
	return nil
}

// validateCourseTitle checks that user input for course title is valid. Returns error type.
func validateCourseTitle(courseTitle string) error {
	if courseTitle != "" {
		if len(courseTitle) > 45 {
			return errors.New(">> Course Title cannot be greater than 45 characters")
		}
	} else {
		return errors.New(">> Course Title cannot be blank")
	}
	return nil
}
