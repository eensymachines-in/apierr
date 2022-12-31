package apierr

/* This was developed as library for error transfer from the various layers of the api
APIs in the HTTP setup typically need an identification so as to determine the corresponding HTTP status code
Also errors are then logged on the server to get the complete debug information
date 			: 31/12/2022
author			: kneerunjun@gmail.com
maintainedby	: eensymachines.in
*/

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/*============================
Error codes as enum
=============================*/

type ErrorCode uint

const (
	MarshallErr = iota + uint(1000)
	DBConnErr
	DBQueryErr
	NilResultErr
	InvldParamErr
	RangeOvrflwErr
)

/*============================
Implementation of custom error
=============================*/
// APIErr : ineternal api error that has the http status code too
type APIErr struct {
	code     ErrorCode  // this is the error code that would translate to http error status
	loginfo  *log.Entry // this shall be the deep information of why the error has occured
	msg      string     // this is the error message
	ctx      string     // this is the location where the error has ocured
	internal error      // internal error so as to trace in logging
}

// Throw : client uses this to make a new error object and throw it over the error interface
// instead of creating an custom object pointer this just encapsulates the error creation
//
/*
func Sample() error{
	return Throw(fmt.Errorf("")).Code(MarshallErr).Message("Invalid request payload").Context("sample").LogInfo(log.WithFields{
		"vartodebug": value,
	})
}
*/
func Throw(e error) *APIErr {
	return &APIErr{
		internal: e,
	}
}

// Error : hence we can send as over error interface
func (apie *APIErr) Error() string {
	return apie.msg
}

// --------------------------
// Below are the chained functions that help set the various props
// --------------------------

// Code : helps to set the code of the Error object
// code is the HTTP code that is used over the response
func (apie *APIErr) Code(code ErrorCode) *APIErr {
	apie.code = code
	return apie
}

// Message : sets the message of the error
func (apie *APIErr) Message(m string) *APIErr {
	apie.msg = m
	return apie
}

// Context : sets the context of the error - generally this is where the error has emanated from
func (apie *APIErr) Context(c string) *APIErr {
	apie.ctx = c
	return apie
}

// LogInfo : sets the logging info and can be used in a chain
func (apie *APIErr) LogInfo(li *log.Entry) *APIErr {
	apie.loginfo = li
	return apie
}

// ------------------------
// this helps in the client side code to get error details
// ------------------------
// toHttpStatus : maps the the code from
func toHttpStatus(code ErrorCode) int {
	codeMap := map[ErrorCode]int{
		ErrorCode(MarshallErr):    http.StatusBadRequest,
		ErrorCode(DBConnErr):      http.StatusBadGateway,
		ErrorCode(DBQueryErr):     http.StatusInternalServerError,
		ErrorCode(RangeOvrflwErr): http.StatusInternalServerError,
		ErrorCode(InvldParamErr):  http.StatusBadRequest,
		ErrorCode(NilResultErr):   http.StatusNotFound,
	}
	httpEr, ok := codeMap[code]
	if !ok {
		return http.StatusInternalServerError
	}
	return httpEr
}

// ToHTTPContext : this converts the error to http context
// From within the handler function (has the gin.Context) this can log and pack it into HttpStatusResponse
//
/*
func Hndler(c *gin.Context) {
	err := Throw(fmt.Errorf("")).Code(MarshallErr).Message("Invalid request payload").Context("sample").LogInfo(log.WithFields{
	"vartodebug": value,
	})
	err.Log().ToHTTPContext(c)
	return
}
*/
func (apie *APIErr) ToHTTPContext(c *gin.Context) {
	if c != nil {
		c.AbortWithStatusJSON(toHttpStatus(apie.code), gin.H{
			"err": apie.Error(),
		})
	}
}

// Log : this persists the error to the log as configured
// incase the log entry is nil, this will print only the error
func (apie *APIErr) Log() *APIErr {
	logStr := fmt.Sprintf("%s: %d: %s", apie.ctx, apie.code, apie.Error())
	if apie.loginfo == nil {
		log.Errorf(logStr)
	} else {
		apie.loginfo.Error(logStr)
	}
	return apie
}
