package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAgent(t *testing.T) {
	pollInterval, reportInterval := time.Duration(2*time.Second), time.Duration(10*time.Second)
	serverAddr := `:8080`

	a := NewAgent(pollInterval, reportInterval, serverAddr)

	assert.Equal(t, pollInterval, a.PollInterval)
	assert.Equal(t, reportInterval, a.ReportInterval)
	assert.Equal(t, serverAddr, a.ServerAddr)
}
