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

		assert.False(t, cfg.Verbose)
		assert.Equal(t, "test", cfg.Environment)
	})

	t.Run("Load_Should_Show_Env_Vars_Values", func(t *testing.T) {
		// Setup
		cfg := &Config{}
		os.Setenv(config.AppEnvironmentKey, "example")
		os.Setenv(config.VerboseKey, "false")
		os.Setenv("EXAMPLE_ENV_AWS_KEY", "aws_key_env")
		os.Setenv("EXAMPLE_ENV_AWS_SECRET", "aws_secret_env")

		defer os.Unsetenv(config.AppEnvironmentKey)
		defer os.Unsetenv(config.VerboseKey)
		defer os.Unsetenv("EXAMPLE_ENV_AWS_KEY")
		defer os.Unsetenv("EXAMPLE_ENV_AWS_SECRET")

		err := cfg.Load(MockConfigPath)
		assert.NoError(t, err, "expected no error from Load")

		// Assertions
		assert.Equal(t, "aws_key_env", cfg.AWS.Key)
		assert.Equal(t, "aws_secret_env", cfg.AWS.Secret)
		assert.False(t, cfg.Verbose)
		assert.Equal(t, "example", cfg.Environment)
	})
}
