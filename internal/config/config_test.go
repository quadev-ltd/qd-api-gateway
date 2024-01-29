package config

import (
	"os"
	"testing"

	"github.com/quadev-ltd/qd-common/pkg/config"
	pkgConfig "github.com/quadev-ltd/qd-common/pkg/config"
	"github.com/stretchr/testify/assert"
)

const (
	MockConfigPath = "./"
)

func TestLoad(t *testing.T) {
	t.Run("Load_Should_Show_File_Values_If_No_Env_Vars", func(t *testing.T) {
		// Setup
		cfg := &Config{}
		os.Setenv(pkgConfig.AppEnvironmentKey, "test")
		os.Setenv(pkgConfig.VerboseKey, "false")

		defer os.Unsetenv(pkgConfig.AppEnvironmentKey)
		defer os.Unsetenv(pkgConfig.VerboseKey)

		err := cfg.Load(MockConfigPath)
		assert.NoError(t, err, "expected no error from Load")

		// Assertions

		assert.Equal(t, true, cfg.TLSEnabled)
		assert.Equal(t, "qd.api.gateway", cfg.GRPC.Host)
		assert.Equal(t, "9000", cfg.GRPC.Port)
		assert.Equal(t, "qd.email.api", cfg.EmailService.Host)
		assert.Equal(t, "9091", cfg.EmailService.Port)
		assert.Equal(t, "qd.visualization.api", cfg.VisualizationService.Host)
		assert.Equal(t, "9092", cfg.VisualizationService.Port)
		assert.Equal(t, "qd.authentication.api", cfg.AuthenticationService.Host)
		assert.Equal(t, "9090", cfg.AuthenticationService.Port)

		assert.False(t, cfg.Verbose)
		assert.Equal(t, "test", cfg.Environment)
	})

	t.Run("Load_Should_Show_Env_Vars_Values", func(t *testing.T) {
		// Setup
		cfg := &Config{}
		os.Setenv(config.AppEnvironmentKey, "test")
		os.Setenv(config.VerboseKey, "false")
		os.Setenv("TEST_ENV_GRPC_HOST", "localhost_env")
		os.Setenv("TEST_ENV_GRPC_PORT", "3333_env")
		os.Setenv("TEST_ENV_EMAIL_SERVICE_HOST", "localhost_env")
		os.Setenv("TEST_ENV_EMAIL_SERVICE_PORT", "3333_env")
		os.Setenv("TEST_ENV_AUTHENTICATION_SERVICE_HOST", "localhost_env")
		os.Setenv("TEST_ENV_AUTHENTICATION_SERVICE_PORT", "4444_env")
		os.Setenv("TEST_ENV_VISUALIZATION_SERVICE_HOST", "localhost_env")
		os.Setenv("TEST_ENV_VISUALIZATION_SERVICE_PORT", "5555_env")

		defer os.Unsetenv(config.AppEnvironmentKey)
		defer os.Unsetenv(config.VerboseKey)
		defer os.Unsetenv("TEST_ENV_GRPC_HOST")
		defer os.Unsetenv("TEST_ENV_GRPC_PORT")
		defer os.Unsetenv("TEST_ENV_EMAIL_SERVICE_HOST")
		defer os.Unsetenv("TEST_ENV_EMAIL_SERVICE_PORT")

		err := cfg.Load(MockConfigPath)
		assert.NoError(t, err, "expected no error from Load")

		// Assertions
		assert.Equal(t, "localhost_env", cfg.GRPC.Host)
		assert.Equal(t, "3333_env", cfg.GRPC.Port)
		assert.Equal(t, "localhost_env", cfg.EmailService.Host)
		assert.Equal(t, "3333_env", cfg.EmailService.Port)
		assert.Equal(t, "localhost_env", cfg.AuthenticationService.Host)
		assert.Equal(t, "4444_env", cfg.AuthenticationService.Port)
		assert.Equal(t, "localhost_env", cfg.VisualizationService.Host)
		assert.Equal(t, "5555_env", cfg.VisualizationService.Port)

		assert.False(t, cfg.Verbose)
		assert.Equal(t, "test", cfg.Environment)
	})
}
