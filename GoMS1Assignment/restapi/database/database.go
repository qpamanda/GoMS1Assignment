// Package database implements the connection to the database server at the designated port.
// It performs the CRUD operations to the database as invoked by the REST API.
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// courseInfo struct for the json
type courseInfo struct {
	Title string `json:"Title"`
}

// Config struct to maintain DB configuration properties
type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

// Connector variable used for CRUD operations
var DB *sql.DB

// Connect creates the database connection
func Connect(connectionString string) error {
	var err error
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	} else {
		fmt.Println("Database opened")
	}
	return nil
}

// GetConnectionString formats the database connection string and returns it.
func GetConnectionString(config Config) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.User, config.Password, config.ServerName, config.DB)
	return connectionString
}

// AddCourse implements the sql operations to insert a new course as invoked by the REST API.
func AddCourse(courseID string, courseTitle string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Courses (CourseID, CourseTitle, Created_DT, LastModified_DT) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(courseID, courseTitle, time.Now(), time.Now())
	if err != nil {
		panic("error executing sql insert")
	}
}

// UpdateCourse implements the sql operations to update a course as invoked by the REST API.
func UpdateCourse(courseID string, courseTitle string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Courses SET CourseTitle=?, LastModified_DT=? WHERE CourseID=?")
	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(courseTitle, time.Now(), courseID)
	if err != nil {
		panic("error executing sql update")
	}
}

// DeleteCourse implements the sql operations to delete a course as invoked by the REST API.
func DeleteCourse(courseID string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("DELETE FROM Courses WHERE CourseID=?")
	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(courseID)
	if err != nil {
		panic("error executing sql update")
	}
}

// GetCourse implements the sql operations to retrieve a course as invoked by the REST API.
func GetCourse(courseID string) map[string]courseInfo {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var courseTitle string

	// Instantiate courses
	var courses = make(map[string]courseInfo)

	query := "SELECT CourseID, CourseTitle FROM Courses WHERE CourseID=?"

	results, err := DB.Query(query, courseID)
	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			err := results.Scan(&courseID, &courseTitle)
			if err != nil {
				panic("error getting results from sql select")
			}
			courses[courseID] = courseInfo{courseTitle}
		}
		return courses
	}
}

// GetAllCourses implements the sql operations to retrieve all courses as invoked by the REST API.
func GetAllCourses() map[string]courseInfo {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var courseID string
	var courseTitle string

	// Instantiate courses
	var courses = make(map[string]courseInfo)

	query := "SELECT CourseID, CourseTitle FROM Courses"

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	} else {
		for results.Next() {
			err := results.Scan(&courseID, &courseTitle)
			if err != nil {
				panic("error getting results from sql select")
			}

			courses[courseID] = courseInfo{courseTitle}
		}
		return courses
	}
}

// ValidKey implements the sql operations to retrieve the access key and validate if the
// key provided is valid. Returns a bool.
func ValidKey(key string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	query := "SELECT KeyID FROM APIKeys WHERE KeyID=?"

	results, err := DB.Query(query, key)
	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			return true
		}
	}
	return false
}
