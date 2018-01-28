// TODO?
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// HttpStatus represents the an http status code.
type HttpStatus int

// TODO: comment, refactor, and clean up code
// TODO: credit others (licenses) where credit is due

// TODO: enable caching via middleware?
// TODO: vendor dependencies
// TODO: credit sources (gin, google/uuid, etc.)
// TODO: add a root "/" path that lists APIs (to make services discoverable)?
// TODO: test AJAX support with your webpage (to update the default hours if
// different)

// TODO: enable accounts for updating hours, ~(uploading pictures)?, etc.
// TODO: OAuth / OpenId ... ?

// TODO: add API for retrieving gallery images?

// TODO: finish API and WRITE TESTS!!!!!!!!!
// TODO: commit your changes to the git repository

// TODO: add error logging
// TODO: add support for error / log tracking

// TODO: write up a swagger doc
// TODO: general HTML/other format from swagger doc and include in webpage...

// TODO: separate out endpoints into separate services (for organization and
// efficiency) (does gin take care of some of this for you?) (would this over
// complicate things because you'd need to setup a proxy / router perhaps to
// route urls to different microservices?

// TODO: replace files with a backing database? (e.g. MySQL) or at least
// implement the database version and save the other and you can go back

// HTTPS security middleware options
var httpsSecurityOptions = secure.Options{
	// TODO: REMOVE (only for development)
	IsDevelopment: true,

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

// TODO: vendor (and document?) all of your dependencies
// TODO: generate any documentation you might want
// TODO: build and deploy your server
// TODO: host server somewhere (droplet?)
// - TODO: (optional) purchase a hostname
// TODO: add webpage prototype to your blog -> document what you did, your
// goals, design decisions, takeaways, topics for improvement and further
// exploration (e.g. making a game called krepl!)
// TODO: contact Mr. Canamucio to get his permision to use the webpage in your
// portfolio (with citation)
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

	// secure authentication middleware
	// TODO: FIX THIS -> YOUR USING IT WRONG
	//authFunc := func(c *gin.Context) {
	//restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static,
	//restgate.Config{
	//Key:    []string{"12345"},
	//Secret: []string{"secret"},
	//})
	//}

	// middleware
	// TODO: REMOVE?
	router.Use(func(c *gin.Context) { c.Writer.Header().Set("Access-Control-Allow-Origin", "*") })
	router.Use(secureFunc)

	router.GET("/hours", GetHours)
	router.GET("/hours/findActive", GetActiveHours)
	router.GET("/hours/findByName/:name", GetHoursByName)
	router.PUT("/hours", PutHours)
	router.DELETE("/hours/findByName/:name", DeleteHours)

	// TODO: REMOVE
	//authorized.Use(AuthRequired())
	//{
	//authorized.PUT("/hours", PutHours)
	//authorized.DELETE("/hours/findByName/:name", DeleteHours)
	//}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}
