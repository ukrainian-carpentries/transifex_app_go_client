package transifex_app_go_client

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// The function configures the API server internal logger
func (t *TransifexAppClient) configureLogger(config *Config) error {

	// Configuration of the logger destination: stdout or a file
	logFile, err := setLoggerDestination(t.l, config.LogDestination)
	if err != nil {
		return err
	}

	// Configuration of the logger formatter: text or json
	if err := setLoggerFormatter(t.l, config.LogFormatter); err != nil {
		return err
	}

	// Configuration of the logger level: Trace, Debug, Info, Warning, Error, Fatal or Panic
	if err := setLoggerLevel(t.l, config.LogLevel); err != nil {
		return err
	}

	// Configure the function, that will be called in the case of "Ctrl+C".
	// The function closes the log file and exits.
	setupCloseHandler(t.l, logFile)

	t.l.Debugf("logger level is set to '%s'", config.LogLevel)
	t.l.Debugf("logger destination is set to '%s'", config.LogDestination)
	t.l.Debugf("logger formatter is set to '%s'", config.LogFormatter)
	return nil
}

func setLoggerDestination(logger *logrus.Logger, dst string) (*os.File, error) {
	// Set the os.Stdout or a file for writing the log messages
	if len(dst) == 0 || strings.ToLower(dst) == "stdout" {

		// If the destination is not configured or set to stdout explicitly
		logger.SetOutput(os.Stdout)
		return nil, nil
	}

	// Open a file for the logger output
	logFile, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("logger: New(): unable to open the file '%s' for writing: %s", dst, err.Error())
	}

	// Redirect the logger output to the file
	logger.SetOutput(logFile)

	return logFile, nil
}

func setLoggerLevel(logger *logrus.Logger, levelStr string) error {

	// If the logging level is not configured,the "info" logging level is used,
	// since an http.Server and httputil.ReverseProxy use it when send
	// messages to a given Writer.
	if levelStr == "" {
		level, err := logrus.ParseLevel("info")
		if err != nil {
			return fmt.Errorf("logger: New(): unable to set the logging level 'info': %s", err.Error())
		}
		logger.SetLevel(level)
		return nil
	}

	// Set the logging level
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return fmt.Errorf("logger: New(): unable to set the logging level '%s': %s", levelStr, err.Error())
	}
	logger.SetLevel(level)

	return nil
}

// The function sets the logger formatter (mainly logrus.TextFormatter{} or logrus.JSONFormatter{})
func setLoggerFormatter(logger *logrus.Logger, formatter string) error {

	// Set the logger formatter
	switch strings.ToLower(formatter) {

	// If not configured, the JSON formatter is used as the default one
	case "":
		fallthrough
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})

	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})

	default:
		return fmt.Errorf("logger: New(): unknown logger formatter '%s'", formatter)
	}
	return nil
}

func setupCloseHandler(logger *logrus.Logger, f *os.File) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Debug("- 'Ctrl + C' was pressed in the Terminal. Terminating...")

		// Close a log file before exit
		if f != nil {
			f.Close()
		}
		os.Exit(0)
	}()
}
