package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/eensymachines-in/apierr"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestErrLog: this will try to make and throw error also will try to test the logging
func TestErrLog(t *testing.T) {
	// Setting up logging
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	// By default the log output is stdout and the level is info
	log.SetOutput(os.Stdout)     // FLogF will set it main, but dfault is stdout
	log.SetLevel(log.DebugLevel) // default level info debug but FVerbose will set it main

	err := apierr.Throw(fmt.Errorf("sample error for test"))
	assert.NotNil(t, err, "unexpected nil value of error from Throw")
	err.Code(apierr.ErrorCode(apierr.NilResultErr)).Context("testing/sample-err").LogInfo(log.WithFields(log.Fields{
		"sample": "sample value",
	})).Message("sample message from within the testing package").Log()
}

func TestToHTTPCode(t *testing.T) {
	code := apierr.ToHttpStatus(apierr.ErrorCode(apierr.InvldParamErr))
	assert.Equal(t, 400, code, "unexpected http error code")
	code = apierr.ToHttpStatus(apierr.ErrorCode(apierr.DBConnErr))
	assert.Equal(t, 502, code, "unexpected http error code")
}
