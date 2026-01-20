package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_ValidConfig(t *testing.T) {
	content := `
temporal:
  host: localhost:7233
  task_queue: TEST_QUEUE

nats:
  url: nats://localhost:4222
  timeout: 3s

activity:
  timeout: 10s
`
	tmpFile := filepath.Join(t.TempDir(), "config.yaml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	assert.Equal(t, "localhost:7233", cfg.Temporal.Host)
	assert.Equal(t, "TEST_QUEUE", cfg.Temporal.TaskQueue)
	assert.Equal(t, "nats://localhost:4222", cfg.Nats.URL)
	assert.Equal(t, 3*time.Second, cfg.Nats.Timeout)
	assert.Equal(t, 10*time.Second, cfg.Activity.Timeout)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "invalid.yaml")
	err := os.WriteFile(tmpFile, []byte("invalid: yaml: content:"), 0644)
	require.NoError(t, err)

	_, err = Load(tmpFile)
	assert.Error(t, err)
}

func TestConfig_EnvOverrides(t *testing.T) {
	content := `
temporal:
  host: localhost:7233
  task_queue: TEST_QUEUE

nats:
  url: nats://localhost:4222
  timeout: 2s

activity:
  timeout: 5s
`
	tmpFile := filepath.Join(t.TempDir(), "config.yaml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	// 设置环境变量
	os.Setenv("TEMPORAL_HOST", "temporal.example.com:7233")
	os.Setenv("NATS_URL", "nats://nats.example.com:4222")
	defer func() {
		os.Unsetenv("TEMPORAL_HOST")
		os.Unsetenv("NATS_URL")
	}()

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	assert.Equal(t, "temporal.example.com:7233", cfg.Temporal.Host)
	assert.Equal(t, "nats://nats.example.com:4222", cfg.Nats.URL)
}
