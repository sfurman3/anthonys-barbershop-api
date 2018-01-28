package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// HttpStatus represents the an http status code.
type HttpStatus int

// HTTPS security middleware options
var httpsSecurityOptions = secure.Options{
	AllowedHosts:          []string{"sethfurman.github.io"},
	HostsProxyHeaders:     []string{},
	SSLRedirect:           false,
	SSLHost:               "",
	SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	STSSeconds:            315360000,
	STSIncludeSubdomains:  true,
	STSPreload:            true,
	FrameDeny:             true,
	ContentTypeNosniff:    true,
	BrowserXssFilter:      true,
	ContentSecurityPolicy: "script-src $NONCE",
	PublicKey:             `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`,
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// HTTPS security middlware
	secureMiddleware := secure.New(httpsSecurityOptions)
	secureFunc := func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			log.Println(err)
			c.Abort()
			return
		}
	}

	// middleware
	//// TODO: REMOVE?
	//router.Use(func(c *gin.Context) { c.Writer.Header().Set("Access-Control-Allow-Origin", "*") })
	router.Use(secureFunc)

	router.GET("/hours", GetHours)
	router.GET("/hours/findActive", GetActiveHours)
	router.GET("/hours/findByName/:name", GetHoursByName)
	//router.PUT("/hours", PutHours)
	//router.DELETE("/hours/findByName/:name", DeleteHours)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}
