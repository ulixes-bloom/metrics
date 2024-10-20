package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/service"
	"github.com/ulixes-bloom/ya-metrics/internal/agent/storage"
)

func TestNewAgent(t *testing.T) {
	pollInterval, reportInterval := time.Duration(2*time.Second), time.Duration(10*time.Second)
	serverAddr := `:8080`

	m := storage.NewMemStorage()
	s := service.NewService(m)
	c := NewClient(s, pollInterval, reportInterval, serverAddr)

	assert.Equal(t, pollInterval, c.PollInterval)
	assert.Equal(t, reportInterval, c.ReportInterval)
	assert.Equal(t, serverAddr, c.ServerAddr)
}
