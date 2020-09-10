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
	ErrCode int    // this is the error code that would translate to http error status
	Msg     string // this is the error message
	Context string // this is the location where the error has ocured
}

// APIErrrFrom : creates a new instance of the error from other lib error
func APIErrrFrom(e error, context string) *APIErr {
	// This will default to internal error
	return &APIErr{500, e.Error(), context}
}

// NewAPIErrr : creates a new error instance of the APIErr
func NewAPIErrr(code int, msg, context string) *APIErr {
	return &APIErr{code, msg, context}
}
func (apie *APIErr) Error() string {
	return apie.Msg
}

// ToHTTPContext : this converts the error to http context
func (apie *APIErr) ToHTTPContext(c *gin.Context) {
	c.AbortWithStatusJSON(apie.ErrCode, gin.H{
		"err": apie.Error(),
	})
	return
}

// Log : this persists the error to the log as configured
func (apie *APIErr) Log() {
	log.Errorf("%s: %d: %s", apie.Context, apie.ErrCode, apie.Error())
}
