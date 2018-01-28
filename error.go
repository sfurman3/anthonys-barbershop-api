package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Error represents the content of an API error.
type ErrorContent struct {
	Id     uuid.UUID `json:"id"`
	Status int       `json:"status"`
	Detail string    `json:"detail"`
	URL    string    `json:"url"`
}

// Error represents an API error object for this service.
type Error struct {
	Err ErrorContent `json:"error"`
}

// NewError creates a new Error object.
func NewError(c *gin.Context, status int, detail string) Error {
	return Error{ErrorContent{Id: uuid.New(), URL: c.Request.URL.Path, Status: status, Detail: detail}}
}
