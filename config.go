package main

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const binaryName = "hasap-alerts-handler"

// Config new config from yaml
func Config(flagSet *flag.FlagSet) (*viper.Viper, error) {
	config := viper.New()

	err := config.BindPFlags(flagSet)
	if err != nil {
		return nil, errors.Wrap(err, "could not bind config to CLI flags")
	}

	// try to get the "config" value from the bound "config" CLI flag
	path := config.GetString("config")
	if path != "" {
		// try to manually load the configuration from the given path
		err = loadConfigurationFromFile(config, path)
	} else {
		// otherwise try viper's auto-discovery
		err = loadConfigurationAutomatically(config)
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not load configuration file")
	}

	setLogLevel(config.GetString("log-level"))

	return config, nil
}

func loadConfigurationAutomatically(config *viper.Viper) error {
	config.SetConfigName(binaryName)
	config.AddConfigPath("/usr/etc/")

	err := config.ReadInConfig()
	if err == nil {
		log.Info("Using config file: ", config.ConfigFileUsed())
		return nil
	}

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Infof("Could not discover configuration file: %s", err)
		log.Info("Fallingback to default values.")
		log.Warn("NO Alerts will be sent to alertmanager in case of failures. See alertmanagerIP variable for this")
		return nil
	}

	return errors.Wrap(err, "could not load automatically discovered config file")
}

// loads configuration from an explicit file path
func loadConfigurationFromFile(config *viper.Viper, path string) error {
	// we hard-code the config type to yaml, otherwise ReadConfig will not load the values
	// see https://github.com/spf13/viper/issues/316
	config.SetConfigType("yaml")

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not open file")
	}
	defer file.Close()

	err = config.ReadConfig(file)
	if err != nil {
		return errors.Wrap(err, "could not read file")
	}

	log.Info("Using config file: ", path)

	return nil
}

func setLogLevel(level string) {
	switch level {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Warnln("Unrecognized minimum log level; using 'info' as default")
	}
}
