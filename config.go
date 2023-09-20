package transifex_app_go_client

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	// ErrInvalidLogLevel        = errors.New("invalid logging level")
	// ErrInvalidLogFormatter    = errors.New("invalid logging formatter")
	// ErrInvalidLogDestination  = errors.New("invalid logging destination")
	// ErrInvalidApiUrl          = errors.New("invalid transifex url")
	// ErrNoToken                = errors.New("no access token was provided")
	ErrNoConfigFilePath = errors.New("no config file path was provided")

// ErrUnableToOpenConfigFile = errors.New("unabe to open the config file")
)

// The Config struct stores main Transifex API Cient configuration parameters
type Config struct {
	LogLevel       string `yaml:"log_level"`
	LogDestination string `yaml:"log_destination"`
	LogFormatter   string `yaml:"log_formatter"`
	// ApiURL         string `yaml:"api_url"`
	// Token          string `yaml:"api_token"`
}

// The function creates a config with the default parameter values,
// then overrides the parameter values with ones from the input file.
// If any of the config parameters has invalid value, the method returns a coresponding error
func NewConfigFromFile(path string) (*Config, error) {
	config := &Config{
		LogLevel:       "error",
		LogDestination: "stdout",
		LogFormatter:   "text",
		// ApiURL:         "https://rest.api.transifex.com",
		// Token:          "",
	}

	// Override the default parameter values with the values from the input file
	err := config.updateFromFile(path)
	if err != nil {
		return nil, err
	}

	// return the config and the result of its checking as the error message
	return config, config.check()
}

// The function parses a configuration yaml file and overrides
// the corresponding parameter values in the config struct
func (c *Config) updateFromFile(configPath string) error {

	// If the config file path was not provided
	if configPath == "" {
		return ErrNoConfigFilePath
	}

	// Open the config file
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("unable to open the configuration file '%s': %s", configPath, err.Error())
	}
	defer file.Close()

	// Init a new yaml decoder
	d := yaml.NewDecoder(file)

	// Decode the configuration from the config file
	err = d.Decode(&c)
	if err != nil {
		return fmt.Errorf("unable to decode the configuration file '%s' as yaml: %s", configPath, err.Error())
	}

	return nil
}

// The function checks, whether all mandatory parameters of the config have values
func (c *Config) check() error {
	// var fields string = ""

	// if c.Token == "" {
	// 	fields += "api_token,"
	// }

	// if fields != "" {
	// 	return fmt.Errorf("the following parameters of the config file have empty values: '%s'", strings.TrimSuffix(fields, ","))
	// }

	return nil
}
