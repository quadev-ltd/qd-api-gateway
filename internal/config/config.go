package config

import (
	"fmt"

	commonAWS "github.com/quadev-ltd/qd-common/pkg/aws"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	"github.com/rs/zerolog/log"
)

// Config is the configuration of the application
type Config struct {
	Verbose     bool
	Environment string
	AWS         commonAWS.Config
	TLSEnabled  bool
}

// Load loads the configuration from the given path yml file
func (config *Config) Load(path string) error {
	env := commonConfig.GetEnvironment()
	config.Environment = env
	config.Verbose = commonConfig.GetVerbose()

	log.Info().Msgf("Loading configuration for environment: %s", env)
	vip, err := commonConfig.SetupConfig(path, env)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}
	if err := vip.Unmarshal(&config); err != nil {
		return fmt.Errorf("Error unmarshaling configuration: %v", err)
	}

	return nil
}
