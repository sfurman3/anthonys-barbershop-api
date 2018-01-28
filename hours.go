package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

// HoursDir is the directory containing all business hours.
const HoursDir = "hours"

// HoursBaseURL is the base URL of the business hours API endpoint.
const HoursBaseURL = "/hours"

// Extension of hours files.
const HoursFileExtension = ".json"

// Error messages
const (
	NOT_FOUND_FMT      = `The requested URL "%s" was not found on this server.`
	INT_SERVR_FMT      = `The requested URL "%s" generated an internal server error.`
	JSON_UNMARSHAL_FMT = `The requested URL "%s" generated an JSON unmarshaling internal server error.`
	NO_HOURS_FOUND_MSG = `No hours were found on the server.`
)

// GenericTime represents a generic time (hours and minutes).
type GenericTime struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

// GenericHours represents generic business hours for a single day (hours and
// minutes to hours and minutes)
//
// NOTE: If IsClosed is true, then the business is considered closed on that
// day.
type GenericHours struct {
	Weekday   time.Weekday `json:"weekday"`
	StartTime GenericTime  `json:"start_time"`
	EndTime   GenericTime  `json:"end_time"`
	IsClosed  bool         `json:"is_closed"`
}

// SpecificHours represents business hours on a specific date or sequece of
// dates (hours and minutes).
//
// NOTE: If the start time and end time are on different days, then the hours
// are treated as a sequence of dates with the same hours (i.e.
// StartTime.Date() to EndTime.Date() from StartTime.Clock() to
// EndTime.Clock()).
//
// NOTE: If IsClosed is true, then the business is considered closed on the
// given day or sequence of days.
type SpecificHours struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsClosed  bool      `json:"is_closed"`
}

// HoursSet represents a set of 0+ GenericHours and 0+ SpecificHours.
//
// Title indicates the user-facing title of the set of hours.
// Active indicates whether the hours should be visible to users.
type HoursSet struct {
	Title    string          `json:"title"`
	Active   bool            `json:"active"`
	Generic  []GenericHours  `json:"generic_hours"`
	Specific []SpecificHours `json:"specific_hours"`
}

// PUTHoursSet is the type used for posting a new set of hours.
//
// Name is a unique identifier for the set and is, for instance, when GETing an
// hours set.
//
// NOTE: using an existing name for an hours set will overwrite those hours.
type PUTHoursSet struct {
	Name string   `json:"name"`
	Set  HoursSet `json:"hours_set"`
}

// HoursNamesResponse represents a list of all names, one for each set of
// hours.
type HoursNamesResponse struct {
	HoursNames []string `json:"hours_names"`
}

// hoursFileName returns a string of the file name (path included) for the
// given hours set name.
func hoursFileName(name string) string {
	return HoursDir + string(os.PathSeparator) + name + HoursFileExtension
}

// getHours returns a list of names of all available sets of hours and an
// HttpStatus code indicating success or failure.
//
// NOTE: Creates the "hours" directory if it doesn't already exist.
func getHours(c *gin.Context) ([]string, HttpStatus) {
	hoursDirFile, err := os.Open(HoursDir)
	if os.IsNotExist(err) {
		// create HoursDir
		rw_ := 06
		r__ := 04
		userPerms := rw_ << 6
		groupPerms := r__ << 3
		otherPerms := r__
		os.Mkdir(HoursDir, os.FileMode(userPerms|groupPerms|otherPerms))

		// no hours were found
		return nil, http.StatusNotFound
	} else if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}

	readAllNames := 0
	hoursFileNames, err := hoursDirFile.Readdirnames(readAllNames)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError
	}
	if len(hoursFileNames) == 0 {
		log.Println("No hours")
		return nil, http.StatusNotFound
	}

	var hoursNames []string
	for _, fname := range hoursFileNames {
		if extension := path.Ext(fname); extension == HoursFileExtension {
			hoursName := fname[:len(fname)-len(extension)]
			hoursNames = append(hoursNames, hoursName)
		}
	}

	return hoursNames, http.StatusOK
}

// GetHours responds with a list of all names of hours currently stored on
// the server.
func GetHours(c *gin.Context) {
	hoursNames, httpStatus := getHours(c)
	status := int(httpStatus)
	switch status {
	case http.StatusNotFound:
		c.JSON(status, NewError(c, status, NO_HOURS_FOUND_MSG))
	case http.StatusOK:
		c.JSON(http.StatusOK, HoursNamesResponse{hoursNames})
	default:
		detail := fmt.Sprintf(INT_SERVR_FMT, c.Request.URL.Path)
		c.JSON(status, NewError(c, status, detail))
	}
}

func hoursAreActive(hoursName string) (bool, HttpStatus) {
	rawJSON, err := ioutil.ReadFile(hoursFileName(hoursName))
	if os.IsNotExist(err) {
		return false, http.StatusNotFound
	} else if err != nil {
		log.Println(err)
		return false, http.StatusInternalServerError
	}

	var hours HoursSet
	err = json.Unmarshal(rawJSON, &hours)
	if err != nil {
		log.Println(err)
		return false, http.StatusInternalServerError
	}

	return hours.Active, http.StatusOK
}

func GetActiveHours(c *gin.Context) {
	hoursNames, httpStatus := getHours(c)
	status := int(httpStatus)
	if status != http.StatusOK {
		switch status {
		case http.StatusNotFound:
			c.JSON(status, NewError(c, status, NO_HOURS_FOUND_MSG))
		default:
			c.JSON(status, NewError(c, status,
				fmt.Sprintf(INT_SERVR_FMT, c.Request.URL.Path)))
		}
		return
	}

	var activeHours []string
	for _, name := range hoursNames {
		active, httpStatus := hoursAreActive(name)
		status := int(httpStatus)
		if status != http.StatusOK {
			// NOTE: this is an internal server error even if the
			// response is http.StatusNotFound because the hours
			// should exist since we just read the name with
			// getHours().
			c.JSON(status, NewError(c, status,
				fmt.Sprintf(INT_SERVR_FMT, c.Request.URL.Path)))
			return
		}

		if active {
			activeHours = append(activeHours, name)
		}
	}

	c.JSON(http.StatusOK, HoursNamesResponse{activeHours})
}

// GetHoursByName responds with the business hours associated with the given hours
// name (either regular business hours or special hours (e.g. holiday hours)).
func GetHoursByName(c *gin.Context) {
	name := c.Param("name")
	rawJSON, err := ioutil.ReadFile(hoursFileName(name))
	if os.IsNotExist(err) {
		status := http.StatusNotFound
		detail := fmt.Sprintf(NOT_FOUND_FMT, c.Request.URL.Path)
		c.JSON(status, NewError(c, status, detail))
		return
	} else if err != nil {
		log.Println(err)
		status := http.StatusInternalServerError
		detail := fmt.Sprintf(INT_SERVR_FMT, c.Request.URL.Path)
		c.JSON(status, NewError(c, status, detail))
		return
	}

	var businessHours HoursSet
	err = json.Unmarshal(rawJSON, &businessHours)
	if err != nil {
		log.Println(err)
		status := http.StatusInternalServerError
		detail := fmt.Sprintf(JSON_UNMARSHAL_FMT, c.Request.URL.Path)
		c.JSON(status, NewError(c, status, detail))
		return
	}

	c.JSON(http.StatusOK, businessHours)
}

// PutHours adds or overwrites a requested hours set.
func PutHours(c *gin.Context) {
	var hours PUTHoursSet
	err := c.ShouldBindJSON(&hours)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	data, err := json.Marshal(hours.Set)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// determine whether to return 201: Created or 200: OK
	fname := hoursFileName(hours.Name)
	status := http.StatusOK
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		// path/to/whatever does not exist
		status = http.StatusCreated
	}

	rw_ := 06
	r__ := 04
	userPerms := rw_ << 6
	groupPerms := r__ << 3
	otherPerms := r__
	err = ioutil.WriteFile(fname, data, os.FileMode(userPerms|groupPerms|otherPerms))
	if err != nil {
		log.Println(err)
		status := http.StatusInternalServerError
		detail := fmt.Sprintf(INT_SERVR_FMT, c.Request.URL.Path)
		c.JSON(status, NewError(c, status, detail))
		return
	}

	c.Writer.Header().Set("Location", HoursBaseURL+"/"+"findByName/"+hours.Name)
	c.Writer.WriteHeader(status)
}

// DeleteHours deletes a requested hours set.
func DeleteHours(c *gin.Context) {
	name := c.Param("name")

	err := os.Remove(hoursFileName(name))
	if err != nil {
		log.Println(err)
		if os.IsNotExist(err) {
			status := http.StatusNotFound
			detail := fmt.Sprintf(NOT_FOUND_FMT, c.Request.URL.Path)
			c.JSON(status, NewError(c, status, detail))
			return
		} else {
			status := http.StatusInternalServerError
			detail := fmt.Sprintf(INT_SERVR_FMT, c.Request.URL.Path)
			c.JSON(status, NewError(c, status, detail))
			return
		}
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
