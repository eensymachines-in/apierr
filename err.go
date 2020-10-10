package apierr

/* This package is simply a sugar coating on the error in any of the microservices
Any function now can send not only the error message but also a error code that is compatible with http codess
*/
import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// APIErr : ineternal api error that has the http status code too
type APIErr struct {
	code    int    // this is the error code that would translate to http error status
	msg     string // this is the error message
	ctx     string // this is the location where the error has ocured
	loginfo string // this shall be the deep information of why the error has occured
}

// APIErrrFrom : creates a new instance of the error from other lib error
func APIErrrFrom(e error, context, log string) *APIErr {
	// This will default to internal error
	return &APIErr{500, e.Error(), context, log}
}

// NewAPIErrr : creates a new error instance of the APIErr
func NewAPIErrr(code int) *APIErr {
	result := &APIErr{}
	result.code = code
	return result
}
func (apie *APIErr) Error() string {
	return apie.msg
}

// --------------------------
// Below are the chained functions that help set the various props
// --------------------------

// Code : helps to set the code of the Error object
// code is the HTTP code that is used over the response
func (apie *APIErr) Code(code int) *APIErr {
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
func (apie *APIErr) LogInfo(li string) *APIErr {
	apie.loginfo = li
	return apie
}

// ------------------------
// this helps in the client side code to get error details
// ------------------------

// ToHTTPContext : this converts the error to http context
func (apie *APIErr) ToHTTPContext(c *gin.Context) {
	c.AbortWithStatusJSON(apie.code, gin.H{
		"err": apie.Error(),
	})
	return
}

// Log : this persists the error to the log as configured
func (apie *APIErr) Log() {
	log.Errorf("%s: %d: %s", apie.ctx, apie.code, apie.Error())
}
