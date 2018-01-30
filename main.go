package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// HttpStatus represents the an http status code.
type HttpStatus int

// ApisResponse represents the response containing all available endpoints,
// returned from the server's root.
type ApisResponse struct{
	Apis []string `json:"apis"`
}

func main() {
	// set mode to 'release'
	gin.SetMode(gin.ReleaseMode)

	// set up logging
	logName := fmt.Sprintf("/root/logs/anthonys-barbershop-api/log%d.log", time.Now().Unix())
	logFile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	ginLogName := fmt.Sprintf("/root/logs/anthonys-barbershop-api/gin-log%d.log", time.Now().Unix())
	ginLogFile, err := os.OpenFile(ginLogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer ginLogFile.Close()
	gin.DefaultWriter = io.MultiWriter(ginLogFile)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// middleware
	router.Use(func(c *gin.Context) { c.Writer.Header().Set("Access-Control-Allow-Origin", "*") })

	router.GET("/", func (c *gin.Context) {
		c.JSON(http.StatusOK, ApisResponse{
			Apis: []string{"/hours"},
		})
	})
	router.GET("/hours", GetHours)
	router.GET("/hours/findActive", GetActiveHours)
	router.GET("/hours/findByName/:name", GetHoursByName)
	//router.PUT("/hours", PutHours)
	//router.DELETE("/hours/findByName/:name", DeleteHours)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}
